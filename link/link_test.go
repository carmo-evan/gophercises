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

func assert(t *testing.T, value, expected interface{}) {
	if value != expected {
		t.Errorf("returned %#v. Expected %#v", value, expected)
	}
}

/*TestLinkParse tests the Parse function of the link package. It checks that the length of the returned
array is correct, as well as the accuracy of the text and the href properties of the Links */

func TestLinkParseEx1(t *testing.T) {

	file, err := os.Open("./ex1.html")
	if err != nil {
		t.Error(err)
	}
	r := bufio.NewReader(file)
	links := link.Parse(r)

	assert(t, len(links), 1)
}
func TestLinkParseEx2(t *testing.T) {

	file, err := os.Open("./ex2.html")
	if err != nil {
		t.Error(err)
	}
	r := bufio.NewReader(file)
	links := link.Parse(r)

	assert(t, len(links), 2)
	assert(t, links[0].Href, "https://www.twitter.com/joncalhoun")
	assert(t, links[0].Text, "Check me out on twitter")
	assert(t, links[1].Href, "https://github.com/gophercises")
	assert(t, links[1].Text, "Gophercises is on Github !")
}

func TestLinkParseEx3(t *testing.T) {

	file, err := os.Open("./ex3.html")
	if err != nil {
		t.Error(err)
	}
	r := bufio.NewReader(file)
	links := link.Parse(r)

	assert(t, len(links), 3)
	assert(t, links[0].Href, "#")
	assert(t, links[0].Text, "Login")
	assert(t, links[1].Href, "/lost")
	assert(t, links[1].Text, "Lost? Need help?")
	assert(t, links[2].Href, "https://twitter.com/marcusolsson")
	assert(t, links[2].Text, "@marcusolsson")
}

func TestLinkParseEx4(t *testing.T) {
	file, err := os.Open("./ex4.html")
	if err != nil {
		t.Error(err)
	}
	r := bufio.NewReader(file)
	links := link.Parse(r)

	assert(t, len(links), 1)

	assert(t, links[0].Href, "/dog-cat")

	assert(t, links[0].Text, "dog cat")
}

func TestLinkParseEx5(t *testing.T) {
	file, err := os.Open("./ex5.html")
	if err != nil {
		t.Error(err)
	}
	r := bufio.NewReader(file)
	links := link.Parse(r)

	assert(t, len(links), 1)

	assert(t, links[0].Href, "/attributes")

	assert(t, links[0].Text, "Testing multiple attrs")
}
