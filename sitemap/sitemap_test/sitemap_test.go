package sitemap_test

import (
	"../mapbuilder"
	"encoding/xml"
	"log"
	"testing"
)

func TestMapBuilder(t *testing.T) {

	siteMap, err := mapbuilder.GetMap(`gophercises.com`)
	if err != nil {
		t.Error(err)
	}
	var urlset mapbuilder.Urlset
	if err := xml.Unmarshal([]byte(siteMap), &urlset); err != nil {
		t.Error(err)
	}
	log.Print(siteMap)
}
