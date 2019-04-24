package link_test

import (
	"./link"
	"bufio"
	"os"
	"testing"
)

func handleErr(e error) {
	if e != nil {
		panic(e)
	}
}

/*TestLinkParse tests the Parse function of the link package. It checks that the length of the returned
array is correct, as well as the accuracy of the text and the href properties of the Links */
func TestLinkParseEx4(t *testing.T) {
	// var links []Link
	file, err := os.Open("./ex4.html")
	if err != nil {
		t.Error(err)
	}
	r := bufio.NewReader(file)
	links := link.Parse(r)
	if len(links) != 1 {
		t.Errorf("returned %d links. Should have returned 1", len(links))
	}
	if links[0].Href != "/dog-cat" {
		t.Errorf("Captured %v as the href. Should have captured '/dog-cat'", links[0].Href)

	}
	if links[0].Text != "dog cat " {
		t.Errorf("Captured %v as the text. Should have captured 'dog cat '", links[0].Text)
	}
}
