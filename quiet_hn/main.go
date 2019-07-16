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

		sort.Ints(ids)
		ids = ids[:numStories]

		tasks := make(chan int, len(ids))
		stories := stories{}
		workersCount := numStories
		wg := sync.WaitGroup{}

		wg.Add(workersCount)

		for i := 0; i < workersCount; i++ {
			go func() {
				defer wg.Done()
				for id := range tasks {
					hnItem, _ := client.GetItem(id) // reach out to API to get full item details
					item := parseHNItem(hnItem)
					stories.mu.Lock()
					stories.items = append(stories.items, item)
					stories.mu.Unlock()
				}
			}()
		}

		for _, id := range ids {
			tasks <- id
		}

		close(tasks)

		wg.Wait()

		sort.Slice(stories.items, func(i, j int) bool {
			return stories.items[i].ID < stories.items[j].ID
		})

		data := templateData{
			Stories: stories.items,
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
	items []item
	mu    sync.RWMutex
}

type templateData struct {
	Stories []item
	Time    time.Duration
}
