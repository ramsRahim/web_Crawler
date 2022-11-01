package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	log "github.com/llimllib/loglevel"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type Link struct {
	url   string
	text  string
	depth int
}

type HttpError struct {
	original string
}

func (self Link) String() string {
	spacer := strings.Repeat("\t", self.depth)
	return fmt.Sprintf("%s%s (%d) - %s", spacer, self.text, self.depth, self.url)
}

func (self Link) Valid() bool {
	if self.depth >= MaxDepth {
		return false
	}

	if len(self.text) == 0 {
		return false
	}
	if len(self.url) == 0 || strings.Contains(strings.ToLower(self.url), "javascript") {
		return false
	}

	return true
}

func (self HttpError) Error() string {
	return self.original
}

var MaxDepth = 2

func LinkReader(resp *http.Response, depth int) []Link {
	page := html.NewTokenizer(resp.Body)
	links := []Link{}

	var start *html.Token
	var text string

	for {
		_ = page.Next()
		token := page.Token()
		if token.Type == html.ErrorToken {
			break
		}

		if start != nil && token.Type == html.TextToken {
			text = fmt.Sprintf("%s%s", text, token.Data)
		}

		if token.DataAtom == atom.A {
			switch token.Type {
			case html.StartTagToken:
				if len(token.Attr) > 0 {
					start = &token
				}
			case html.EndTagToken:
				if start == nil {
					log.Warnf("Link End found without Start: %s", text)
					continue
				}
				link := NewLink(*start, text, depth)
				if link.Valid() {
					links = append(links, link)
					log.Debugf("Link Found %v", link)
				}

				start = nil
				text = ""
			}
		}
	}

	log.Debug(links)
	return links
}

func HtmlParse(resp *http.Response) string {
	page := html.NewTokenizer(resp.Body)

	//var start *html.Token
	text := ""

	for {
		_ = page.Next()
		token := page.Token()
		if token.Type == html.ErrorToken {
			break
		}

		/* if start != nil && token.Type == html.TextToken {
			text = fmt.Sprintf("%s%s", text, token.Data)
		} */

		if token.DataAtom == atom.P {
			text += token.Data
			fmt.Printf("%s", text)
		}
		if token.DataAtom == atom.Span {
			text += token.Data
			fmt.Printf("%s", text)
		}
	}

	return text
}

func NewLink(tag html.Token, text string, depth int) Link {
	link := Link{text: strings.TrimSpace(text), depth: depth}

	for i := range tag.Attr {
		if tag.Attr[i].Key == "href" {
			link.url = strings.TrimSpace(tag.Attr[i].Val)
		}
	}
	return link
}

func recurDownloader(url string, depth int) string {
	page, err := downloader(url)
	if err != nil {
		log.Error(err)
		return "error"
	}
	text := ""
	text += HtmlParse(page)
	//fmt.Printf("%s", text)
	links := LinkReader(page, depth)
	for _, link := range links {
		if depth+1 < MaxDepth {
			recurDownloader(link.url, depth+1)
		}
	}

	return text
}

func downloader(url string) (resp *http.Response, err error) {
	log.Debugf("Downloading %s", url)
	resp, err = http.Get(url)
	if err != nil {
		log.Debugf("Error: %s", err)
		return
	}

	if resp.StatusCode > 299 {
		err = HttpError{fmt.Sprintf("Error (%d): %s", resp.StatusCode, url)}
		log.Debug(err)
		return
	}
	return

}

func main() {

	log.SetPriorityString("info")
	log.SetPrefix("crawler")

	log.Debug(os.Args)

	if len(os.Args) < 2 {
		log.Fatalln("Missing Url arg")
	}

	file_data := recurDownloader(os.Args[1], 0)

	f, err := os.Create("output.txt")
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()
	n, err := f.WriteString(file_data)
	if err != nil {
		panic(err)
	}
	fmt.Printf("wrote %d bytes\n", n)
	f.Sync()

}
