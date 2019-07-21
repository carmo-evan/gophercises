package handlers

import (
	"fmt"
	"net/http"
)

func PanicDemo(w http.ResponseWriter, r *http.Request) {
	funcThatPanics()
}

func PanicAfterDemo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Hello!</h1>")
	funcThatPanics()
}

func funcThatPanics() {
	panic("Oh no!")
}

func Hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "<h1>Hello!</h1>")
}
