package models

import (
	"errors"
	"strconv"
)

type Metric struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func (m *Metric) SetValue(v string) error {
	switch m.MType {
	case "gauge":
		value, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return err
		}
		m.Value = &value
	case "counter":
		value, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return err
		}
		m.Delta = &value
	default:
		return errors.New("unknown metric type")
	}

	return nil
}

func (m *Metric) GetValue() (string, error) {
	switch m.MType {
	case "gauge":
		return strconv.FormatFloat(*m.Value, 'f', -1, 64), nil
	case "counter":
		return strconv.FormatInt(*m.Delta, 10), nil
	}

	return "", errors.New("unknown metric type")
}
