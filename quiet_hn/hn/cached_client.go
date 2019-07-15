package hn

import (
	"strconv"
	"time"

	"github.com/patrickmn/go-cache"
)

type cachedClient struct {
	client Client
	cache  *cache.Cache
}

// NewCachedClient initializes a cachedClient instance, as well as the
// underlying live client instance, as well as an in-memory cache.
func NewCachedClient() Client {
	c := &cachedClient{}
	c.client = &client{}
	c.cache = cache.New(5*time.Minute, 10*time.Minute)
	return c
}

// TopItems returns the ids of roughly 450 top items in decreasing order. These
// should map directly to the top 450 things you would see on HN if you visited
// their site and kept going to the next page.
//
// TopItmes does not filter out job listings or anything else, as the type of
// each item is unknown without further API calls.
// This implementation leverages caching, if possible, or delegates to the live
// implementation.
func (c *cachedClient) TopItems() ([]int, error) {
	cached, found := c.cache.Get("topItems")
	if found {
		return cached.([]int), nil
	}
	// call the underlying client
	live, err := c.client.TopItems()
	c.cache.Set("topItems", live, cache.DefaultExpiration)
	return live, err
}

// GetItem will return the Item defined by the provided ID.
func (c *cachedClient) GetItem(id int) (Item, error) {
	strID := strconv.Itoa(id)
	cached, found := c.cache.Get(strID)
	if found {
		return cached.(Item), nil
	}
	item, err := c.client.GetItem(id)

	if err != nil {
		return item, err
	}

	c.cache.Set(strID, item, cache.DefaultExpiration)

	return item, nil
}
