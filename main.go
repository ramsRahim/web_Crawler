package main

import (
	"flag"
	"os"

	"github.com/cobaltspeech/log"
	"github.com/ramsRahim/webCrawer/http"
)

func main() {
	logger := log.NewLeveledLogger()

	flagURL := flag.String("url", "http://feeds.bbci.co.uk/news/england/london/rss.xml", "url from which to parse text")
	isXml := flag.Bool("XML", true, "checks if the link is xml")
	flagOutput := flag.String("output", "crawler_output.txt", "name of the output file")
	url := *flagURL
	logger.Debug("msg", "reading url", "url", url)

	var data []byte
	var err error

	if *isXml {
		links, err := http.ReadRSS(url)
		if err != nil {
			logger.Error("msg", "unable to read url", "url", url, "error", err)
		}
		for _, link := range links {
			text, err := http.GetText(link)
			if err != nil {
				logger.Error("msg", "unable to read url", "url", url, "error", err)
			}
			data = append(data, text...)
		}
	} else {
		data, err = http.GetText(url)
		if err != nil {
			logger.Error("msg", "unable to read url", "url", url, "error", err)
		}
	}

	f, err := os.OpenFile(*flagOutput, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		logger.Error("msg", "unable to write file", "file", f, "error", err)
	}
	if n, err := f.WriteString(string(data)); err != nil {
		logger.Error("msg", "unable to write bytes", "bytes", n, "error", err)
	}

}
