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
		ids, err := client.TopItems() // returns around 450 unique ids
		if err != nil {
			http.Error(w, "Failed to load top stories", http.StatusInternalServerError)
			return
		}
		s := newStories() //holds map of items, slice of ids, mutex
		wg := &sync.WaitGroup{}
		for _, id := range ids {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				hnItem, err := client.GetItem(id) // reach out to API to get full item details
				if err != nil {
					return
				}
				item := parseHNItem(hnItem)
				if isStoryLink(item) { // if it's the type we want, let's save it
					s.mu.Lock()
					defer s.mu.Unlock()
					s.items[id] = item        // save item
					s.ids = append(s.ids, id) // save id
					//issue here is, we are still reaching out to the API 450 times
					// Yes, we do it concurrently, but that's a lot of times.
					// Benchmark ends up being the same as hitting it 30 times sequentially.
					// so concurrency right now does not represent any performance gain
				}
			}(id)

		}
		wg.Wait()
		var items []item
		sort.Ints(s.ids)
		for _, id := range s.ids {
			item, _ := s.items[id]
			items = append(items, item)
			if len(items) == numStories { //this is where we are capping the displayed items to 30
				break
			}
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
