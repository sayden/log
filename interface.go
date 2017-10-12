package log

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// Interface represents the API of both Logger and Entry.
type Interface interface {
	WithFields(fields Fielder) Interface
	WithField(key string, value interface{}) Interface
	WithError(err error) Interface
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
	Fatal(msg string)
	Debugf(msg string, v ...interface{})
	Infof(msg string, v ...interface{})
	Warnf(msg string, v ...interface{})
	Errorf(msg string, v ...interface{})
	Fatalf(msg string, v ...interface{})
	Trace(msg string) Interface
	addons
	telemetryAddons
}

type telemetryAddons interface {
	Inc(string, float64) Interface
}

type addons interface {
	SetMessage(string)
	setStart(time time.Time)
	GetLevel() Level
	GetFields() Fields
	GetMessage() string
	finalize(level Level, msg string) Interface
	mergedFields() Fields
	GetTimestamp() time.Time
	Stop(err *error)
	GetTelemetry() Telemetry
}

type Telemetry interface {
	WithTags(...string) Interface
	Inc(name string, value float64) Interface

	//Prometheus
	SetPrometheusInc(name string, counter *prometheus.CounterVec)

	//StatsD
	SetNamespace(s string)
}
