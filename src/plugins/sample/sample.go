package main

import (
	"math/rand"

	"../../types"
)

var Resolution int64 = 10000

func GetMetrics() (m []*types.Metric) {
	var metric *types.Metric

	m = []*types.Metric{}
	metric = new(types.Metric)

	metric.Name = "test.metric"
	metric.Labels = map[string]string{
		"foo":  "bar",
		"fizz": "buzz",
	}
	metric.Type = "gauge"
	metric.Help = "just a test metric"
	metric.Value = rand.Int()
	m = append(m, metric)

	return m
}
