package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync"
	"time"

	"gophercises/quiet_hn/hn"
)

func main() {
	// parse flags
	var port, numStories int
	flag.IntVar(&port, "port", 8080, "the port to start the web server on")
	flag.IntVar(&numStories, "num_stories", 30, "the number of top stories to display")
	flag.Parse()

	tpl := template.Must(template.ParseFiles("./index.gohtml"))
	client := hn.NewCachedClient()

	http.HandleFunc("/", handler(client, numStories, tpl))

	// Start the server
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func handler(client hn.Client, numStories int, tpl *template.Template) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ids, err := client.TopItems() // returns around 500 unique ids
		if err != nil {
			http.Error(w, "Failed to load top stories", http.StatusInternalServerError)
			return
		}
		tasks := make(chan int, len(ids))
		results := make(chan item, len(ids))
		workersCount := 30
		wg := sync.WaitGroup{}
		wg.Add(workersCount)
		for i := 0; i < workersCount; i++ {
			go func() {
				defer wg.Done()
				for id := range tasks {
					hnItem, err := client.GetItem(id) // reach out to API to get full item details
					if err != nil {
						panic(err)
					}
					item := parseHNItem(hnItem)
					results <- item
				}
			}()
		}

		for _, id := range ids {
			tasks <- id
		}

		close(tasks)

		wg.Wait()

		var items []item

		for item := range results {
			if isStoryLink(item) {
				items = append(items, item)
			}
		}

		close(results)

		sort.Slice(items, func(i, j int) bool {
			return items[i].ID < items[j].ID
		})

		data := templateData{
			Stories: items[:30],
			Time:    time.Now().Sub(start),
		}
		err = tpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Failed to process the template", http.StatusInternalServerError)
			return
		}
	})
}

func isStoryLink(item item) bool {
	return item.Type == "story" && item.URL != ""
}

func parseHNItem(hnItem hn.Item) item {
	ret := item{Item: hnItem}
	url, err := url.Parse(ret.URL)
	if err == nil {
		ret.Host = strings.TrimPrefix(url.Hostname(), "www.")
	}
	return ret
}

// item is the same as the hn.Item, but adds the Host field
type item struct {
	hn.Item
	Host string
}

type stories struct {
	items map[int]item
	ids   []int
	mu    sync.RWMutex
}

func newStories() *stories {
	var s stories
	s.items = make(map[int]item)
	return &s
}

type templateData struct {
	Stories []item
	Time    time.Duration
}
