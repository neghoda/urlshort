package urlshort

import (
	"encoding/json"
	"net/http"

	"gopkg.in/yaml.v2"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if target, ok := pathsToUrls[r.URL.Path]; ok {
			http.Redirect(w, r, target, http.StatusMovedPermanently)
		}
		fallback.ServeHTTP(w, r)
	}
}

type pathURL struct {
	Path string `yaml:"path" json:"path"`
	URL  string `yaml:"url" json:"url"`
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
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var parsedYAML []pathURL
	if err := yaml.Unmarshal(yml, &parsedYAML); err != nil {
		return nil, err
	}

	pathsToUrls := make(map[string]string)
	for _, v := range parsedYAML {
		pathsToUrls[v.Path] = v.URL
	}
	return MapHandler(pathsToUrls, fallback), nil
}

// JSONHandler will parse the provided JSON and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the JSON, then the
// fallback http.Handler will be called instead.
//
// JSON is expected to be in the format:
//		[{
//     	 "path": "/some-path"
//       "url": "https://www.some-url.com/demo"
//		}]
// The only errors that can be returned all related to having
// invalid JSON data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func JSONHandler(jsn []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var parsedJSON []pathURL
	if err := json.Unmarshal(jsn, &parsedJSON); err != nil {
		return nil, err
	}

	pathsToUrls := make(map[string]string)
	for _, v := range parsedJSON {
		pathsToUrls[v.Path] = v.URL
	}
	return MapHandler(pathsToUrls, fallback), nil
}