package gobacktest

import (
	"errors"
)

var (
	// ErrInvalidKey describes an error for an invalid key parameter.
	ErrInvalidKey = errors.New("invalid key")
)

// Metric holds metric propertys to a data point.
type Metric map[string]float64

// Add a value to the metrics map.
func (m Metric) Add(key string, value float64) error {
	if m == nil {
		m = make(map[string]float64)
	}

	if key == "" {
		return ErrInvalidKey
	}

	m[key] = value
	return nil
}

// Get returns a metric by name, if not found it returns false.
func (m Metric) Get(key string) (float64, bool) {
	value, ok := m[key]
	return value, ok
}
