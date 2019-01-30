package promhttp_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/travelaudience/go-promhttp"
)

func TestServeMux(t *testing.T) {
	mux := &promhttp.ServeMux{
		ServeMux: &http.ServeMux{},
	}
	reg := prometheus.NewRegistry()

	mux.Handle("/a", http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		t.Log("handler a")
	}))
	mux.Handle("/b", http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		t.Log("handler b")
	}))
	mux.HandleFunc("/c", func(_ http.ResponseWriter, _ *http.Request) {
		t.Log("handler c")
	})
	reg.Register(mux)

	srv := httptest.NewServer(mux)
	defer srv.Close()

	for i := 0; i < 10; i++ {
		if _, err := http.Get(srv.URL + "/a"); err != nil {
			t.Fatalf("unexpected error: %s", err.Error())
		}
		if _, err := http.Get(srv.URL + "/b"); err != nil {
			t.Fatalf("unexpected error: %s", err.Error())
		}
		if _, err := http.Get(srv.URL + "/c"); err != nil {
			t.Fatalf("unexpected error: %s", err.Error())
		}
	}

	assertMetrics(t, reg, 15)
}
