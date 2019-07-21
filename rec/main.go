package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"strconv"
	"strings"
	"gophercises/recover/handlers"
	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
)


func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/panic/", handlers.PanicDemo)
	mux.HandleFunc("/panic-after/", handlers.PanicAfterDemo)
	mux.HandleFunc("/", handlers.Hello)
	os.Setenv("ENV", "DEV") // TODO read from a config fil
	log.Fatal(http.ListenAndServe(":8080", recoverMiddleware(mux, os.Getenv("ENV"))))
}


func recoverMiddleware(next http.Handler, env string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buf := NewResponseBuffer()
		defer func() {
			if r := recover(); r != nil {

				log.Println("recovering from ", r)
				stack := string(debug.Stack())
				filename, line := parseFilenameAndLineFromPanicStack(stack)
				log.Printf("%v", stack)

				if env == "DEV" {
					w.Header().Add("Content-Type", "text/html")
					source := getSourceFromFile(filename)
					printStack(w, stack)
					printSource(w, source, line)
					return
				}

				w.WriteHeader(http.StatusInternalServerError)
				body := "Oops. Something went wrong!"
				w.Write([]byte(body))
				return
			}
			buf.Flush(w)
		}()
		next.ServeHTTP(buf, r)
	})
}

func printSource(w io.Writer, source string, line int) {
	lexer := lexers.Get("go")
	style := styles.Get("monokai")
	ranges := [][2]int{
		{line, line},
	}
	formatter := html.New(html.WithLineNumbers(), html.LineNumbersInTable(), html.HighlightLines(ranges))
	iterator, err := lexer.Tokenise(nil, source)
	if err != nil {
		log.Fatal(err)
	}
	err = formatter.Format(w, style, iterator)
	if err != nil {
		log.Fatal(err)
	}
}

func printStack(w io.Writer, stack string) {
	lexer := lexers.Get("bash")
	style := styles.Get("monokai")
	formatter := html.New(html.Standalone())
	iterator, err := lexer.Tokenise(nil, stack)
	if err != nil {
		log.Fatal(err)
	}
	err = formatter.Format(w, style, iterator)
	if err != nil {
		log.Fatal(err)
	}
}

func getSourceFromFile(filename string) string {
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	return string(dat)

}

func parseFilenameAndLineFromPanicStack(stack string) (filename string, line int) {
	// line 0 is go routine info
	// line 1 (odds) is func info, line 2 (evens) is file info
	// look for runtime/panic, return next file line (+2)
	lines := strings.Split(stack, "\n")
	for i, l := range lines {
		if strings.Contains(l, "runtime/panic") {
			fileLine := lines[i+2]
			fmt.Println(fileLine)
			line := parseLineNumberFromStackLine(fileLine)
			// file name is anything before colon
			colonIndex := strings.Index(fileLine, ":")
			fileName := fileLine[:colonIndex]
			return strings.TrimSpace(fileName), line
		}
	}
	return "", 0
}

func parseLineNumberFromStackLine(stackLine string) (ln int) {
	//line number is anything after colon, before whitespace
	fmt.Println(stackLine)
	colonIndex := strings.Index(stackLine, ":")
	stackLine = stackLine[colonIndex+1:] // remove first part
	whiteSpaceIndex := strings.Index(stackLine, " ")
	if whiteSpaceIndex > 0 {
		stackLine = stackLine[:whiteSpaceIndex]
	}
	ln, err := strconv.Atoi(stackLine)
	if err != nil {
		log.Fatal(err)
	}
	return ln
}
