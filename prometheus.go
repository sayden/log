package log

import (
	"github.com/prometheus/client_golang/prometheus"
)

type prom struct {
	id             string
	tags           []string
	counterMetrics map[string]*prometheus.CounterVec
}

func TelemetryPrometheus() Telemetry {
	return &prom{counterMetrics: make(map[string]*prometheus.CounterVec)}
}

func (p *prom) WithTags(tags ...string) Interface {
	p.tags = tags

	return Log
}

func (p *prom) Inc(name string, value float64) Interface {
	p.counterMetrics[name].WithLabelValues(p.tags...).Add(value)

	return Log
}

func (p *prom) SetPrometheusInc(name string, counter *prometheus.CounterVec) {
	p.counterMetrics[name] = counter
	prometheus.MustRegister(counter)
}

func (p *prom) SetNamespace(s string) {

}
