package main

import (
	"math/rand"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	Resolution int64 = 10000

	RandNumber = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "some_rand_number",
		Help: "Current random number between 0 and 100",
	})
)

func GetCollector() prometheus.Collector {
	return RandNumber
}

func UpdateMetric() {
	RandNumber.Set(float64(rand.Intn(99)))
}
