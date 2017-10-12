package log

import (
	"log"

	statsd_client "github.com/DataDog/datadog-go/statsd"
	"github.com/prometheus/client_golang/prometheus"
)

type statsd struct {
	id   string
	tags []string
	c    *statsd_client.Client
}

func TelemetryStatsD() Telemetry {
	c, err := statsd_client.New("127.0.0.1:8125")
	if err != nil {
		log.Fatal(err)
	}

	return &statsd{
		c: c,
	}
}

func (p *statsd) WithTags(tags ...string) Interface {
	p.tags = tags
	return Log
}

func (p *statsd) Inc(name string, value float64) Interface {
	if err := p.c.Incr(name, nil, value); err != nil {
		Error(err.Error())
	}

	return Log
}

func (p *statsd) SetPrometheusInc(name string, counter *prometheus.CounterVec) {
	// N/A
}

func (p *statsd) SetNamespace(s string) {
	p.c.Namespace = s
}
