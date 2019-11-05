package main

import (
	"fmt"
	"net/http"
	"os"
	"plugin"

	"./types"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
)

type metricsPlugin struct {
	name        string
	resolution  int
	metricsFunc plugin.Symbol
	object      *plugin.Plugin
}

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

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("data"))
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
	p.object = obj
	p.metricsFunc = symbol
	p.resolution = *resolution.(*int)

	return p
}

func main() {
	log.SetLevel(log.DebugLevel)

	var p = loadPlugin("sample.so")
	m := p.metricsFunc.(func() []*types.Metric)()
	fmt.Printf("metric name: %s\n", m[0].Name)

	log.Debug("starting server ...")
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/metrics", handleMetricsEndpoint).Methods("GET")
	http.ListenAndServe(fmt.Sprintf("%s:%s", getenvDefault("LISTEN_ADDRESS", "127.0.0.1"), getenvDefault("LISTEN_PORT", "8080")), router)
}
