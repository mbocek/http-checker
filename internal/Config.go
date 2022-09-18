package internal

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type Metrics struct {
	Port int
}

//goland:noinspection SpellCheckingInspection,SpellCheckingInspection,SpellCheckingInspection,SpellCheckingInspection,SpellCheckingInspection,SpellCheckingInspection
type Check struct {
	Url               string `mapstructure:"url"`
	Timeout           int    `mapstructure:"timeout"`
	CheckPeriod       int    `mapstructure:"check-period"`
	Metric            string `mapstructure:"metric"`
	MetricDescription string `mapstructure:"metric-description"`
	ResponseCode      int    `mapstructure:"response-code"`
}

//goland:noinspection SpellCheckingInspection,SpellCheckingInspection
type Config struct {
	CheckPeriod int     `mapstructure:"check-period"`
	Check       []Check `mapstructure:"check"`
	Metrics     Metrics
}

func ReadConfigFile() Config {
	viper.SetConfigName("http-checker")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/et/http-checker")
	viper.AddConfigPath("config")
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
