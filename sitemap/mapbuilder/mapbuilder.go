package mapbuilder

import (
	"../../link/link"
	"bytes"
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

//Urlset represents the XML structure of a site map
type Urlset struct {
	XMLName xml.Name `xml:"urlset"`
	Text    string   `xml:",chardata"`
	Xmlns   string   `xml:"xmlns,attr"`
	URL     []URL    `xml:"url"`
}

//URL ,not to be confused with net/url.URL, represents the URL field of the xml structure of a site map
type URL struct {
	Text string `xml:",chardata"`
	Loc  string `xml:"loc"`
}

func formatURLString(path string, host string) *url.URL {
	path = strings.TrimRight(path, "/")
	target, _ := url.Parse(path)
	target.Scheme = "https"
	if target.Host == "" && host != "" {
		target.Host = host
	}
	target.Fragment = ""
	return target
}

func getResponseBody(url string) ([]byte, error) {
	target := formatURLString(url, "")
	resp, err := http.Get(target.String())
	if err != nil {
		// handle error
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

//GetMap builds an XML string that can be safely unmarshall into the Urlset type representing
//a sitemap as defined by https://www.sitemaps.org/protocol.html
func GetMap(host string) (string, error) {
	start := time.Now()

	urls := []URL{}
	var urlset Urlset
	seenUrls := &sync.Map{}
	parsePageURLs(host, host, seenUrls, &urls)
	urlset.URL = urls
	urlset.Xmlns = "http://www.sitemaps.org/schemas/sitemap/0.9"
	xmlBytes, err := xml.Marshal(urlset)
	ioutil.WriteFile("./ssitemap.xml", xmlBytes, 0644)

	elapsed := time.Since(start)
	log.Printf("Took %s", elapsed)

	return string(xmlBytes), err
}

func parsePageURLs(target string, host string, seenUrls *sync.Map, urls *[]URL) {
	body, err := getResponseBody(target)
	if err != nil {
		return
	}
	links := link.Parse(bytes.NewReader(body))

	wg := &sync.WaitGroup{}
	routineCount := 0
	for _, v := range links {
		wg.Add(1)
		routineCount++
		go getURLFromLinks(wg, v, seenUrls, host, urls)
	}
	wg.Wait()
}

func getURLFromLinks(wg *sync.WaitGroup, v link.Link, seenUrls *sync.Map, host string, urls *[]URL) {
	defer wg.Done()
	var url URL
	formattedURL := formatURLString(v.Href, host)
	if formattedURL == nil {
		return
	}
	urlstring := formattedURL.String()
	if _, seen := seenUrls.Load(urlstring); seen || formattedURL.Host != host {
		return
	}

	seenUrls.Store(urlstring, true)

	parsePageURLs(urlstring, host, seenUrls, urls)

	url.Loc = urlstring
	*urls = append(*urls, url)
}
