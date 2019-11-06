package main

import (
	"fmt"
	"net/http"
	"os"
	"plugin"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type metricsPlugin struct {
	name          string
	resolution    int64
	collectorFunc func() prometheus.Collector
	metricsFunc   func()
	object        *plugin.Plugin
}

var loadedPlugins []*metricsPlugin

func getenvDefault(name string, defaultValue string) string {
	value := os.Getenv(name)
	if value == "" {
		return defaultValue
	}
	return value
}

func loadPlugin(path string) (p *metricsPlugin) {
	obj, err := plugin.Open(path)
	if err != nil {
		log.Error(err)
		panic(err)
	}
	resolution, err := obj.Lookup("Resolution")
	if err != nil {
		log.Error(err)
		panic(err)
	}
	getCollector, err := obj.Lookup("GetCollector")
	if err != nil {
		log.Error(err)
		panic(err)
	}
	updateMetric, err := obj.Lookup("UpdateMetric")
	if err != nil {
		log.Error(err)
		panic(err)
	}

	p = new(metricsPlugin)
	p.name = path
	p.object = obj
	p.collectorFunc = getCollector.(func() prometheus.Collector)
	p.metricsFunc = updateMetric.(func())
	p.resolution = *resolution.(*int64)

	log.Infof("loaded plugin: %s", p.name)
	return p
}

func runPlugin(p *metricsPlugin) {
	var resolution = time.Duration(p.resolution * int64(time.Millisecond))
	for true == true {
		var start = time.Now()
		p.metricsFunc()
		var sleepTime = resolution - time.Since(start)

		time.Sleep(sleepTime)
		if sleepTime < 0 {
			log.Warnf("Excessive execution time: %s", p.name)
		}
	}
}

func main() {
	log.SetLevel(log.DebugLevel)

	loadedPlugins = make([]*metricsPlugin, 0)

	var p = loadPlugin("sample.so")
	loadedPlugins = append(loadedPlugins, p)
	p = loadPlugin("sample2.so")
	loadedPlugins = append(loadedPlugins, p)

	for _, plugin := range loadedPlugins {
		prometheus.MustRegister(plugin.collectorFunc())
		go runPlugin(plugin)
	}

	log.Debug("starting server ...")
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", getenvDefault("LISTEN_ADDRESS", "127.0.0.1"), getenvDefault("LISTEN_PORT", "8080")), nil))
}
