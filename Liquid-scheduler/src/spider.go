package main

import (
	"net/http"
	"math/rand"
	"time"
	"net/url"
	"strings"
)

type Spider struct {
	UserAgent   string
	Method      string
	URL         string
	ContentType string
	Referer     string
	Data        url.Values
	Response    *http.Response
}

func (spider *Spider) do() error {
	client := &http.Client{}
	req, err := http.NewRequest(spider.Method, spider.URL, strings.NewReader(spider.Data.Encode()))
	if err != nil {
		return err
	}

	if len(spider.ContentType) == 0 {
		spider.ContentType = ""
	}
	req.Header.Set("Content-Type", spider.ContentType)

	/* set user-agent */
	if len(spider.UserAgent) == 0 {
		spider.UserAgent = spider.getUA()
	}
	req.Header.Set("User-Agent", spider.UserAgent)

	if len(spider.Referer) == 0 {
		spider.Referer = ""
	}
	req.Header.Set("Referer", spider.Referer)

	spider.Response, err = client.Do(req)
	if err != nil {
		return err
	}
	return nil
}

func (spider *Spider) getResponse() *http.Response {
	return spider.Response
}

func (spider *Spider) getUA() string {
	rand.Seed(time.Now().Unix())
	UAs := []string{
		"Mozilla/5.0 (X11; Linux i686; rv:64.0) Gecko/20100101 Firefox/64.0",
		"Mozilla/5.0 (Windows NT 6.1; WOW64; rv:64.0) Gecko/20100101 Firefox/64.0",
		"Mozilla/5.0 (X11; Linux i586; rv:63.0) Gecko/20100101 Firefox/63.0",
		"Mozilla/5.0 (Windows NT 6.2; WOW64; rv:63.0) Gecko/20100101 Firefox/63.0",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.14; rv:10.0) Gecko/20100101 Firefox/62.0",
		"Mozilla/5.0 (Macintosh; U; Intel Mac OS X 10.13; ko; rv:1.9.1b2) Gecko/20081201 Firefox/60.0",
		"Mozilla/5.0 (Windows NT 6.1; WOW64; rv:54.0) Gecko/20100101 Firefox/58.0",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.77 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML like Gecko) Chrome/51.0.2704.79 Safari/537.36 Edge/14.14931",
		"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/47.0.2526.111 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_2) AppleWebKit/601.3.9 (KHTML, like Gecko) Version/9.0.2 Safari/601.3.9",
	}
	return UAs[rand.Intn(len(UAs))]
}
