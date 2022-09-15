package main

import (
	"flag"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Metrics struct {
	Port int
}

type Check struct {
	Url               string `mapstructure:"url"`
	Timeout           int    `mapstructure:"timeout"`
	CheckPeriod       int    `mapstructure:"check-period"`
	Metric            string `mapstructure:"metric"`
	MetricDescription string `mapstructure:"metric-description"`
}

type Config struct {
	CheckPeriod int     `mapstructure:"check-period"`
	Check       []Check `mapstructure:"check"`
	Metrics     Metrics
}

func readConfigFile() Config {
	viper.SetConfigName("http-checker")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/et/http-checker")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("reading config file (probably doesn't exists): %w", err))
	}
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		panic(fmt.Errorf("unmarshaling config file: %w", err))
	}
	log.Debug().Interface("Configuration", config).Msg("")
	return config
}

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
	default:
		log.Logger = log.Output(zerolog.New(os.Stdout))
	}
}

func ping(url string, timeout time.Duration) bool {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Error().Err(err).Msg("Error in the request")
		return false
	}
	client := http.Client{Timeout: timeout * time.Second}

	start := time.Now()
	response, err := client.Do(request)
	if err != nil {
		log.Error().Err(err).Msg("Error onß execution of request")
		return false
	}
	log.Info().Int("status code", response.StatusCode).Int64("duration", timeTrack(start).Milliseconds()).Msgf("Execution of %s", url)
	return true
}

func timeTrack(start time.Time) time.Duration {
	return time.Since(start)
}

func recordMetrics(config *Config) {
	for _, check := range config.Check {
		go func(metricCheck Check) {
			metric := promauto.NewGauge(prometheus.GaugeOpts{
				Name: metricCheck.Metric,
				Help: metricCheck.MetricDescription,
			})
			for {
				if ping(metricCheck.Url, time.Duration(metricCheck.Timeout)) {
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
	config := readConfigFile()
	recordMetrics(&config)

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":"+strconv.Itoa(config.Metrics.Port), nil)
}
