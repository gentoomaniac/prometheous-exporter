package main

import (
	"fmt"
	"net/http"
	"os"
	"plugin"
	"strings"
	"time"

	"./types"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
)

type metricsPlugin struct {
	name        string
	resolution  int64
	metricsFunc func() []*types.Metric
	object      *plugin.Plugin
}

var metrics map[string][]*types.Metric

func getenvDefault(name string, defaultValue string) string {
	value := os.Getenv(name)
	if value == "" {
		return defaultValue
	}
	return value
}

// return representation of current state
// needs shared data with a main metrics loop that is periodically collecting metric updates
// https://stackoverflow.com/questions/39207608/how-does-golang-share-variables-between-goroutines
func handleMetricsEndpoint(w http.ResponseWriter, r *http.Request) {
	response := make([]string, 0)

	for _, v := range metrics {
		labelstrings := make([]string, 0)
		for _, e := range v {
			for lk, lv := range e.Labels {
				labelstrings = append(labelstrings, fmt.Sprintf("%s=%s", lk, lv))
			}
			response = append(response, fmt.Sprintf("%s{%s} %d", e.Name, strings.Join(labelstrings, ","), e.Value))
		}
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	for _, line := range response {
		w.Write([]byte(line))
	}
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
	symbol, err := obj.Lookup("GetMetrics")
	if err != nil {
		log.Error(err)
		panic(err)
	}

	p = new(metricsPlugin)
	p.name = path
	p.object = obj
	p.metricsFunc = symbol.(func() []*types.Metric)
	p.resolution = *resolution.(*int64)

	log.Infof("loaded plugin: %s", p.name)
	return p
}

func runPlugin(p *metricsPlugin) {
	var resolution = time.Duration(p.resolution * int64(time.Millisecond))
	for true == true {
		var start = time.Now()

		m := p.metricsFunc()
		metrics[p.name] = m
		var sleepTime = resolution - time.Since(start)

		time.Sleep(sleepTime)
		if sleepTime < 0 {
			log.Warnf("Excessive execution time: %s", p.name)
		}
	}
}

func main() {
	log.SetLevel(log.DebugLevel)

	metrics = make(map[string][]*types.Metric)

	var p = loadPlugin("sample.so")

	go runPlugin(p)

	log.Debug("starting server ...")
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/metrics", handleMetricsEndpoint).Methods("GET")
	http.ListenAndServe(fmt.Sprintf("%s:%s", getenvDefault("LISTEN_ADDRESS", "127.0.0.1"), getenvDefault("LISTEN_PORT", "8080")), router)
}
