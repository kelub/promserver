package stats

import (
	//"github.com/Sirupsen/logrus"
	prom "github.com/prometheus/client_golang/prometheus"
	"time"
)

type PromVec struct {
	namespace string
	subsystem string

	gauge     *prom.GaugeVec
	counter   *prom.CounterVec
	histogram *prom.HistogramVec
}

func NewPromVec(namespace string) *PromVec {
	return &PromVec{
		namespace: namespace,
	}
}

func (p *PromVec) Namespace(namespace string) *PromVec {
	p.namespace = namespace
	return p
}

func (p *PromVec) Subsystem(subsystem string) *PromVec {
	p.subsystem = subsystem
	return p
}

func (p *PromVec) Gauge(name string, labels []string) *PromVec {
	p.gauge = prom.NewGaugeVec(
		prom.GaugeOpts{
			Namespace: p.namespace,
			Subsystem: "",
			Name:      name,
			Help:      name,
		}, labels)
	prom.MustRegister(p.gauge)
	return p
}

func (p *PromVec) Counter(name string, labels []string) *PromVec {
	p.counter = prom.NewCounterVec(
		prom.CounterOpts{
			Namespace: p.namespace,
			Subsystem: "",
			Name:      name,
			Help:      name,
		}, labels)
	prom.MustRegister(p.counter)
	return p
}

func (p *PromVec) Histogram(name string, labels []string, buckets []float64) *PromVec {
	p.histogram = prom.NewHistogramVec(
		prom.HistogramOpts{
			Namespace: p.namespace,
			Subsystem: "",
			Name:      name,
			Help:      name,
			Buckets:   buckets,
		}, labels)
	prom.MustRegister(p.histogram)
	return p
}

func (p *PromVec) Inc(labels []string) {
	if p.counter != nil {
		p.counter.WithLabelValues(labels...).Inc()
	}
	if p.gauge != nil {
		p.gauge.WithLabelValues(labels...).Inc()
	}
}

func (p *PromVec) Dec(labels []string) {
	if p.gauge != nil {
		p.gauge.WithLabelValues(labels...).Dec()
	}
}

func (p *PromVec) Add(labels []string, value float64) {
	if p.counter != nil {
		p.counter.WithLabelValues(labels...).Add(value)
	}
	if p.gauge != nil {
		p.gauge.WithLabelValues(labels...).Add(value)
	}
}

func (p *PromVec) Set(labels []string, value float64) {
	if p.gauge != nil {
		p.gauge.WithLabelValues(labels...).Set(value)
	}
}

func (p *PromVec) HandleTime(labels []string, start time.Time) {
	p.histogram.WithLabelValues(labels...).Observe(time.Since(start).Seconds())
}
