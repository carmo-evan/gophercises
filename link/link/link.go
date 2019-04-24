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

// join concatenates strings
func join(strs ...string) string {
	var sb strings.Builder
	for _, str := range strs {
		sb.WriteString(str)
	}
	return sb.String()
}

//Parse takes a pointer to an HTML file and parses all anchor tags into an array of the Link type
func Parse(r io.Reader) []Link {
	tokenizer := html.NewTokenizer(r)
	depth := 0
	var currentLink Link
	var links []Link
	for {
		tt := tokenizer.Next()
		switch tt {
		case html.ErrorToken:
			err := tokenizer.Err()
			if err == io.EOF {
				return links
			}
		case html.TextToken:
			if depth > 0 {
				text := string(tokenizer.Text())
				currentLink.Text = join(currentLink.Text, text)
			}
		case html.StartTagToken, html.EndTagToken:
			tn, _ := tokenizer.TagName()
			if len(tn) == 1 && tn[0] == 'a' {
				if tt == html.StartTagToken {
					var newLink Link
					currentLink = newLink
					key, val, _ := tokenizer.TagAttr()
					attrName := string(key)
					attrValue := string(val)
					if attrName == "href" {
						currentLink.Href = attrValue
					}
					depth++
				} else {
					if len(tn) == 1 && tn[0] == 'a' {
						links = append(links, currentLink)
					}
					depth--
				}
			}
		}
	}
}
