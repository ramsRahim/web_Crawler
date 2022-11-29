package main

import (
	"flag"
	"os"

	"github.com/cobaltspeech/log"
	"github.com/ramsRahim/webCrawer/http"
)

func main() {
	logger := log.NewLeveledLogger()

	flagURL := flag.String("url", "https://www.cobaltspeech.com/", "url from which to parse text")
	url := *flagURL
	logger.Debug("msg", "reading url", "url", url)
	data, err := http.GetText(url)
	if err != nil {
		logger.Error("msg", "unable to read url", "url", url, "error", err)
	}
	os.Stdout.Write(data)
}
