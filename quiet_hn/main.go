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
		ids, err := client.TopItems()
		if err != nil {
			http.Error(w, "Failed to load top stories", http.StatusInternalServerError)
			return
		}
		s := newStories()
		wg := &sync.WaitGroup{}
		for _, id := range ids {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				hnItem, err := client.GetItem(id)
				if err != nil {
					return
				}
				item := parseHNItem(hnItem)
				if isStoryLink(item) {
					s.mu.Lock()
					defer s.mu.Unlock()
					s.items[id] = item
					s.ids = append(s.ids, id)
					if len(s.items) == numStories {
						// TODO: Cause all other goroutines to quit
						return
					}
				}
			}(id)

		}
		wg.Wait()
		var items []item
		sort.Ints(s.ids)
		for _, id := range s.ids {
			fmt.Println(id)
			item, _ := s.items[id]
			items = append(items, item)
		}
		data := templateData{
			Stories: items,
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
