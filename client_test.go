package promhttp_test

import (
	"crypto/tls"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/travelaudience/go-promhttp"
)

func TestClient_ForRecipient(t *testing.T) {
	r := prometheus.NewRegistry()
	c := promhttp.Client{
		Client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		},
		Registerer: r,
	}

	ca, err := c.ForRecipient("a")
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
	cb, err := c.ForRecipient("b")
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}

	srv := httptest.NewTLSServer(http.NotFoundHandler())
	defer srv.Close()

	url := strings.Replace(srv.URL, "127.0.0.1", "localhost", 1)

	for i := 0; i < 15; i++ {
		if _, err := ca.Get(url); err != nil {
			t.Fatal(err)
		}

		if _, err := cb.Get(url); err != nil {
			t.Fatal(err)
		}
	}

	assertMetrics(t, r, 14)
}
