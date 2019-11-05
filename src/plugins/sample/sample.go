package main

import (
	"../../types"
	log "github.com/Sirupsen/logrus"
)

var Resolution int = 10
var m []*types.Metric

func GetMetrics() []*types.Metric {
	m = []*types.Metric{}
	log.Info("Hello, World\n")
	var metric *types.Metric
	metric = new(types.Metric)
	metric.Name = "test.metric"
	metric.Help = "just a test metric"
	metric.Value = 42
	m = append(m, metric)

	return m
}
