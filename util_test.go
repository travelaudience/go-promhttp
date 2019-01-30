package promhttp_test

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
)

func assertMetrics(t *testing.T, r prometheus.Gatherer, exp int) {
	var count int
	metricFamilies, err := r.Gather()
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
	for _, mf := range metricFamilies {
		count += len(mf.Metric)
		t.Log(mf.GetName(), "-", mf.GetHelp())
	}
	if count != exp {
		t.Errorf("wrong number of metrics, expected %d got %d", exp, count)
	}
}
