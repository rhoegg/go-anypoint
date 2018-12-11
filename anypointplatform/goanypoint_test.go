package anypointplatform_test

import (
	"context"
	"github.com/rhoegg/go-anypoint/anypointplatform"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

var (
	mux *http.ServeMux

	ctx = context.TODO()

	client *anypointplatform.Client

	server *httptest.Server
)

func setup() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	client = anypointplatform.NewClient(nil)
	url, _ := url.Parse(server.URL)
	client.BaseURL = url

}

func teardown() {
	server.Close()
}

func testMethod(t *testing.T, r *http.Request, expected string) {
	if expected != r.Method {
		t.Errorf("Request method = %v, expected %v", r.Method, expected)
	}
}

func handleHttp(t *testing.T, path string, method string, handler func(w http.ResponseWriter, r *http.Request)) {
	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, method)
		handler(w, r)
	})

}
