package promhttp

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
)

func ExampleServeMux() {
	// Create a promhttp ServeMux
	mux := &ServeMux{
		ServeMux: &http.ServeMux{},
	}

	// Attach two endpoints to the mux
	mux.Handle("/path-a", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("path-b handler called")
		w.WriteHeader(http.StatusOK)
	}))
	mux.Handle("/path-b", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("path-b handler called")
	}))

	// Create a test server to easily interact with the newly created mux.
	ts := httptest.NewServer(mux)
	defer ts.Close()

	// Call the endpoint three times.
	for i := 0; i < 3; i++ {
		_, err := http.Get(ts.URL + "/path-a")
		if err != nil {
			log.Fatal(err)
		}
	}

	// The following code is obtaining the Prometheus data. Normally this is
	// done automatically. In this example we are going to get the the counter
	// of how many times an endopint is called along with all of it's labels.

	// Make the chan from which the prom metrics will come in.
	in := make(chan prometheus.Metric)

	// Concurrently fetch the metrics.
	go mux.Collect(in)

	// Create the varible to store the Prom metric
	m := &dto.Metric{}

	// Itterate through the incomming metrics until we get a counter.
	for {
		promMetric := <-in
		promMetric.Write(m)
		if m.Counter != nil {
			break
		}
	}

	// Display the output. First the labels for the metric.
	fmt.Println("Labels:")
	for _, labelPair := range m.Label {
		fmt.Printf("%s -> %s\n", *labelPair.Name, *labelPair.Value)
	}

	// Then the value of the counter.
	fmt.Println("Endpoint calls:", *m.Counter.Value)

	// Output:
	// path-b handler called
	// path-b handler called
	// path-b handler called
	// Labels:
	// code -> 200
	// method -> get
	// path -> /path-a
	// Endpoint calls: 3
}
