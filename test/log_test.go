package test

import (
	"net/http"
	"os"
	"testing"

	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sayden/log"
	"github.com/sayden/log/handlers/text"
)

func Test_Prometheus(t *testing.T) {
	var hdFailures = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "hd_errors_total",
			Help: "Number of hard-disk errors.",
		},
		[]string{"device"},
	)

	log.SetHandler(text.New(os.Stdout))
	log.SetLevel(log.LevelDebug)
	log.SetTelemetry(log.TelemetryPrometheus())
	log.SetPrometheusInc("counter", hdFailures)

	log.Info("hello")

	log.WithTags("a tag").Inc("counter", 1).WithField("a field", "A field value").Info("Warm Hello")

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":8080", nil)
	}()

	time.Sleep(time.Second * 10000)
}

func Test_Statsd(t *testing.T) {
	log.SetHandler(text.New(os.Stdout))
	log.SetLevel(log.LevelDebug)
	log.SetTelemetry(log.TelemetryStatsD())
	log.SetNamespace("mynamespace.")

	log.Info("hello statsd")

	log.WithTags("myapp.myservice").Inc("counter", 1).Info("Warm Hello statsd 1")
	log.WithField("label", "counter").Inc("counter", 1).Info("Warm Hello statsd")
}
