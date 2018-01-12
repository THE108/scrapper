package parser

import (
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gocolly/colly"
)

func visitUrl(fullUrl string, c *colly.Collector, wg *sync.WaitGroup) {
	defer wg.Done()

	hdr := make(http.Header, 2)
	hdr.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) "+
		"Chrome/63.0.3239.84 Safari/537.36")
	hdr.Set("Accept", "*/*")

	if err := c.Request("GET", fullUrl, nil, nil, hdr); err != nil {
		log.Println(err.Error())
	}
}

func Parse(urlList, allowedDomains []string, fn func(*colly.Collector)) {
	c := colly.NewCollector()
	c.AllowedDomains = allowedDomains
	//c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 1})
	c.DisableCookies()
	//c.SetDebugger(debug.Debugger(&debug.LogDebugger{}))

	c.SetRequestTimeout(20 * time.Second)

	c.WithTransport(&http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   20 * time.Second,
		ExpectContinueTimeout: 10 * time.Second,
	})

	fn(c)

	c.OnRequest(func(r *colly.Request) {
		r.Ctx.Put("orig_url", r.URL.String())
	})

	c.OnResponse(func(resp *colly.Response) {
		if resp.StatusCode != http.StatusOK {
			log.Printf("status code: %d headers: %v", resp.StatusCode, *resp.Headers)
		}
	})

	var wg sync.WaitGroup
	wg.Add(len(urlList))
	for _, url := range urlList {
		go visitUrl(url, c, &wg)
	}

	c.Wait()
	wg.Wait()
}
