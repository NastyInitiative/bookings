package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type postData struct {
	key   string
	value string
}

var theTests = []struct {
	name               string
	url                string
	method             string
	params             []postData
	expectedStatusCode int
}{
	{"home", "/", "GET", []postData{}, http.StatusOK},
	{"about", "/about", "GET", []postData{}, http.StatusOK},
	{"generals", "/generals-quarters", "GET", []postData{}, http.StatusOK},
	{"majors", "/majors-suite", "GET", []postData{}, http.StatusOK},
	{"sa", "/search-availability", "GET", []postData{}, http.StatusOK},
	{"contact", "/contact", "GET", []postData{}, http.StatusOK},
	{"mr", "/make-reservation", "GET", []postData{}, http.StatusOK},
	{"post-search-avail", "/search-availability", "POST", []postData{
		{key: "start", value: "01/01/2021"},
		{key: "end", value: "05/01/2021"},
	}, http.StatusOK},
	{"post-search-avail-json", "/search-availability-json", "POST", []postData{
		{key: "start", value: "01/01/2021"},
		{key: "end", value: "05/01/2021"},
	}, http.StatusOK},
	{"make reservation post", "/make-reservation", "POST", []postData{
		{key: "first-name", value: "John"},
		{key: "last-name", value: "Doe"},
		{key: "email", value: "my.me@pipolo.com"},
		{key: "phone", value: "3245556789"},
	}, http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	for _, e := range theTests {
		if e.method == "GET" {
			resp, err := ts.Client().Get(ts.URL + e.url)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}
			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s, expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}
		} else {
			values := url.Values{}
			for _, x := range e.params {
				values.Add(x.key, x.value)
			}
			resp, err := ts.Client().PostForm(ts.URL+e.url, values)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}
			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s, expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}
		}
	}
}
