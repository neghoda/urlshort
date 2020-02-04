package urlshort

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	fallbackBody = "fallback"
	path         = "/example"
	url          = "https://www.google.com"
	fallback     = func(w http.ResponseWriter, r *http.Request) { fmt.Fprint(w, fallbackBody) }
	pathsToUrls  = map[string]string{path: url}
)

func TestMapHandler(t *testing.T) {
	fallbackHandler := http.HandlerFunc(fallback)
	handler := MapHandler(pathsToUrls, fallbackHandler)

	runMapPathTests(t, pathsToUrls, handler)
	runFallbackHandler(t, handler)
}

func TestYAMLHandler(t *testing.T) {
	fallbackHandler := http.HandlerFunc(fallback)
	yaml := fmt.Sprintf(`
  - path: %s
    url: %s
  `, path, url)

	handler, error := YAMLHandler([]byte(yaml), fallbackHandler)
	if error != nil {
		t.Errorf("Error while parsing YAML file - %v", error)
	}
	runMapPathTests(t, pathsToUrls, handler)
	runFallbackHandler(t, handler)
}

func TestJSONHandler(t *testing.T) {
	fallbackHandler := http.HandlerFunc(fallback)
	json := fmt.Sprintf(`[{"path": "%s","url": "%s"}]`, path, url)

	handler, error := JSONHandler([]byte(json), fallbackHandler)
	if error != nil {
		t.Errorf("Error while parsing JSON file - %v", error)
	}
	runMapPathTests(t, pathsToUrls, handler)
	runFallbackHandler(t, handler)
}

func runFallbackHandler(t *testing.T, handler http.HandlerFunc) {
	t.Helper()
	r := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	handler(w, r)
	response := w.Result()

	assertResponseStatus(t, response, http.StatusOK)
	assertResponseBody(t, response, fallbackBody)
}

func runMapPathTests(t *testing.T, pathsToUrls map[string]string, handler http.HandlerFunc) {
	t.Helper()
	for path, expectedLocation := range pathsToUrls {
		r := httptest.NewRequest(http.MethodGet, path, nil)
		w := httptest.NewRecorder()
		handler(w, r)
		response := w.Result()

		assertResponseStatus(t, response, http.StatusMovedPermanently)
		assertResponseLocation(t, response, expectedLocation)

	}
}

func assertResponseStatus(t *testing.T, response *http.Response, expectedStatus int) {
	t.Helper()
	if status := response.StatusCode; status != expectedStatus {
		t.Errorf("Expected %v response status, got -%v", expectedStatus, status)
	}
}

func assertResponseLocation(t *testing.T, response *http.Response, expectedLocation string) {
	t.Helper()
	if location := response.Header.Get("Location"); location != expectedLocation {
		t.Errorf("Expected %v Header Location, got - %v", expectedLocation, location)
	}
}

func assertResponseBody(t *testing.T, response *http.Response, expectedBody string) {
	t.Helper()
	body, _ := ioutil.ReadAll(response.Body)
	if string(body) != expectedBody {
		t.Errorf("Expected %v as fallback body, got - %v", expectedBody, body)
	}
}
