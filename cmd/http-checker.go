package main

import (
	"flag"
	"fmt"
	"github.com/mbocek/http-checker/internal"
	"github.com/mbocek/http-checker/pkg"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func init() {
	const (
		console string = "console"
		json           = "json"
	)
	logFormat := flag.String("log", console, "Log output (console, json)")
	flag.Parse()

	switch *logFormat {
	case console:
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})
	case json:
	default:
		log.Logger = log.Output(zerolog.New(os.Stdout))
	}
}

func recordMetrics(config *internal.Config) {
	for _, check := range config.Check {
		go func(metricCheck internal.Check) {
			metric := promauto.NewGauge(prometheus.GaugeOpts{
				Name: metricCheck.Metric,
				Help: metricCheck.MetricDescription,
			})
			for {
				if pkg.Ping(metricCheck.Url, time.Duration(metricCheck.Timeout), metricCheck.ResponseCode) {
					metric.Set(1)
				} else {
					metric.Set(0)
				}
				time.Sleep(time.Duration(metricCheck.CheckPeriod) * time.Second)
			}
		}(check)
	}
}

func main() {
	config := internal.ReadConfigFile()
	recordMetrics(&config)

	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(":"+strconv.Itoa(config.Metrics.Port), nil)
	if err != nil {
		panic(fmt.Errorf("conot create web server: %w", err))
	}
}
