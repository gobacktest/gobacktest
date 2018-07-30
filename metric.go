package gobacktest

import (
	"errors"
)

// MetricHandler defines the handling of metrics to a data event
type MetricHandler interface {
	Add(string, float64) error
	Get(string) (float64, bool)
}

// Metric holds metric propertys to a data point.
type Metric map[string]float64

// Add ads a value to the metrics map
func (m Metric) Add(key string, value float64) error {
	if m == nil {
		m = make(map[string]float64)
	}

	if key == "" {
		return errors.New("invalid key given")
	}

	m[key] = value
	return nil
}

// Get return a metric by name, if not found it returns false.
func (m Metric) Get(key string) (float64, bool) {
	value, ok := m[key]
	return value, ok
}
