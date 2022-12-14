package http

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/antchfx/xmlquery"
)

func GetText(url string) ([]byte, error) {
	var netClient = &http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := netClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("unable to read response for url: %w", err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	// Find the review items
	doc.Find("body").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the name.
		buf.WriteString(s.Find("p").Text() + "\n")
	})
	return buf.Bytes(), nil
}

func ReadRSS(url string) ([]string, error) {
	var netClient = &http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := netClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("unable to read response for url: %w", err)
	}
	defer resp.Body.Close()

	doc, err := xmlquery.Parse(resp.Body)
	if err != nil {
		return nil, err
	}

	var link []string
	// Find the review items
	for _, n := range xmlquery.Find(doc, "//item/link") {
		link = append(link, n.InnerText())
	}
	return link, nil
}
