package urlshort

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"log"
	"net/http"
)

type yamlConfig struct {
	Path string
	URL  string
}

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mappedURL, ok := pathsToUrls[r.URL.Path]
		if ok {
			http.Redirect(w, r, mappedURL, http.StatusSeeOther)
		} else {
			fallback.ServeHTTP(w, r)
		}
	}
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yaml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedYaml, err := parseYAML(yaml)
	if err != nil {
		return nil, err
	}
	pathMap := buildMap(parsedYaml)
	return MapHandler(pathMap, fallback), nil
}

func parseYAML(data []byte) ([]yamlConfig, error) {

	var parsedYaml []yamlConfig

	err := yaml.Unmarshal([]byte(data), &parsedYaml)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("--- m:\n%v\n\n", parsedYaml)
	return parsedYaml, nil
}

func buildMap(parsedYaml []yamlConfig) map[string]string {
	pathsToUrls := make(map[string]string)
	fmt.Println(parsedYaml)
	for _, pathToURL := range parsedYaml {
		pathsToUrls[pathToURL.Path] = pathToURL.URL
	}
	return pathsToUrls
}
