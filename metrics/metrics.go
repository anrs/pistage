package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/projecteru2/aa/errors"
	"github.com/projecteru2/aa/os"
	"github.com/projecteru2/aa/str"
)

var (
	// DefaultLabels .
	DefaultLabels = []string{"host"}

	// MetricHeartbeatCount .
	MetricHeartbeatCount = "aa_heartbeat_total"
	// MetricErrorCount .
	MetricErrorCount = "aa_error_total"

	metr *Metrics
)

func init() { // nolint
	hn, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	metr = New(hn)
	metr.RegisterCounter(MetricErrorCount, "aa errors", nil)         // nolint
	metr.RegisterCounter(MetricHeartbeatCount, "aa heartbeats", nil) // nolint
}

// Metrics .
type Metrics struct {
	host     string
	counters map[string]*prometheus.CounterVec
	gauges   map[string]*prometheus.GaugeVec
}

// New .
func New(host string) *Metrics {
	return &Metrics{
		host:     host,
		counters: map[string]*prometheus.CounterVec{},
		gauges:   map[string]*prometheus.GaugeVec{},
	}
}

// RegisterCounter .
func (m *Metrics) RegisterCounter(name, desc string, labels []string) error {
	var col = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: name,
			Help: desc,
		},
		str.Merge(labels, DefaultLabels),
	)

	if err := prometheus.Register(col); err != nil {
		return errors.Trace(err)
	}

	m.counters[name] = col

	return nil
}

// RegisterGauge .
func (m *Metrics) RegisterGauge(name, desc string, labels []string) error {
	var col = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: name,
			Help: desc,
		},
		str.Merge(labels, DefaultLabels),
	)

	if err := prometheus.Register(col); err != nil {
		return errors.Trace(err)
	}

	m.gauges[name] = col

	return nil
}

// Incr .
func (m *Metrics) Incr(name string, labels map[string]string) error {
	var col, exists = m.counters[name]
	if !exists {
		return errors.Errorf("collector %s not found", name)
	}

	labels = m.appendLabel(labels, "host", m.host)

	col.With(labels).Inc()

	return nil
}

// Store .
func (m *Metrics) Store(name string, value float64, labels map[string]string) error {
	var col, exists = m.gauges[name]
	if !exists {
		return errors.Errorf("collector %s not found", name)
	}

	labels = m.appendLabel(labels, "host", m.host)

	col.With(labels).Set(value)

	return nil
}

func (m *Metrics) appendLabel(labels map[string]string, key, value string) map[string]string {
	if labels != nil {
		labels[key] = value
	} else {
		labels = map[string]string{key: value}
	}
	return labels
}

// Handler .
func Handler() http.Handler {
	return promhttp.Handler()
}

// IncrError .
func IncrError() {
	Incr(MetricErrorCount, nil) // nolint
}

// IncrHeartbeat .
func IncrHeartbeat() {
	Incr(MetricHeartbeatCount, nil) // nolint
}

// Incr .
func Incr(name string, labels map[string]string) error {
	return metr.Incr(name, labels)
}

// Store .
func Store(name string, value float64, labels map[string]string) error {
	return metr.Store(name, value, labels)
}

// RegisterGauge .
func RegisterGauge(name, desc string, labels []string) error {
	return metr.RegisterGauge(name, desc, labels)
}

// RegisterCounter .
func RegisterCounter(name, desc string, labels []string) error {
	return metr.RegisterCounter(name, desc, labels)
}
