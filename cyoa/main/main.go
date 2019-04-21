package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {

	jsonFile, err := os.Open("../gopher.json")

	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var stories map[string]interface{}

	json.Unmarshal(byteValue, &stories)

	cyoaHandler := getCYOAHandler(stories)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", cyoaHandler)
}

func getCYOAHandler(stories map[string]interface{}) http.HandlerFunc {
	template := template.New("story.html")          // Create a template.
	template, _ = template.ParseFiles("story.html") // Parse template file.

	return func(w http.ResponseWriter, r *http.Request) {
		path := left(r.URL.Path, 1)
		if story, ok := stories[path]; ok {
			template.Execute(w, story)
		} else {
			fmt.Fprint(w, "custom 404")
		}
	}
}

func left(s string, n int) string {
	m := 0
	for i := range s {
		if m >= n {
			return s[i:]
		}
		m++
	}
	return s[:0]
}
