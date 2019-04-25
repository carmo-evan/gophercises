package link

import (
	"golang.org/x/net/html"
	"io"
	"strings"
)

// Link holds two string properties, the text within the original anchor tag and the href attr
type Link struct {
	Href string
	Text string
}

//standardizeSpaces substitutes multiple spaces for just one
func standardizeSpaces(s string) string {
	ss := strings.Fields(s)
	return strings.Join(ss, " ")
}

//assignTextToLink appends the input text to the existing text in the Link
func assignTextToLink(bytes []byte, link *Link) {
	text := standardizeSpaces(string(bytes))
	if len(text) > 0 && len(link.Text) > 0 {
		stringsToJoin := []string{link.Text, text}
		link.Text = strings.Join(stringsToJoin, " ")
	} else if len(text) > 0 {
		link.Text = text
	}
}

func isAnchorLink(tagName []byte) bool {
	return len(tagName) == 1 && tagName[0] == 'a'
}

func assignHrefToLink(link *Link, tokenizer *html.Tokenizer) {
	for {
		key, val, moreAttr := tokenizer.TagAttr()
		attrName := string(key)
		attrValue := string(val)
		if attrName == "href" {
			link.Href = attrValue
			break
		}
		if !moreAttr {
			break
		}
	}
}

//Parse takes a pointer to an HTML file and parses all anchor tags into an array of the Link type
func Parse(r io.Reader) []Link {
	tokenizer := html.NewTokenizer(r)
	depth := 0
	var currentLink Link
	var links []Link
	for {
		tokenType := tokenizer.Next()
		switch tokenType {
		case html.ErrorToken:
			err := tokenizer.Err()
			if err == io.EOF {
				return links
			}
		case html.TextToken:
			if depth > 0 {
				assignTextToLink(tokenizer.Text(), &currentLink)
			}
		case html.StartTagToken:
			tagName, _ := tokenizer.TagName()
			if isAnchorLink(tagName) {
				var newLink Link
				currentLink = newLink
				assignHrefToLink(&currentLink, tokenizer)
				depth++
			}
		case html.EndTagToken:
			tn, _ := tokenizer.TagName()
			if isAnchorLink(tn) {
				links = append(links, currentLink)
				depth--
			}
		}
	}
}
