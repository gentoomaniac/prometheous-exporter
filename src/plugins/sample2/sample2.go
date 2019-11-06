package main

import (
	"math/rand"

	"../../types"
)

var Resolution int64 = 1000

func GetMetrics() (m []*types.Metric) {
	var metric *types.Metric

	m = []*types.Metric{}
	metric = new(types.Metric)

	metric.Name = "fizz.buzz.metric"
	metric.Labels = map[string]string{
		"foobar": "barfoo",
	}
	metric.Type = "gauge"
	metric.Help = "some other metric"
	metric.Value = rand.Intn(99)
	m = append(m, metric)

	return m
}
