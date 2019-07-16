package main

import (
	"gophercises/quiet_hn/hn"
	"sort"
	"sync"
	"testing"
	"time"
)

func BenchmarkQuietHn(b *testing.B) {
	client := hn.NewCachedClient()
	numStories := 30
	ids, err := client.TopItems() // returns around 500 unique ids
	if err != nil {
		b.Error(err)
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

	_ = templateData{
		Stories: stories.items[:30],
		Time:    time.Now().Sub(time.Now()),
	}
}
