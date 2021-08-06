package config

import (
	"github.com/vrischmann/envconfig"
)

type AppConfig struct {
	Addr      string   `envconfig:"default=0.0.0.0:8080"`
	RateLogin int      `envconfig:"default=10"`
	RatePwd   int      `envconfig:"default=100"`
	RateIP    int      `envconfig:"default=10000"`
	Whitelist []string `envconfig:"optional"`
	Blacklist []string `envconfig:"optional"`
	LogLevel  string   `envconfig:"default=DEBUG"`
}

func GetConfig() (AppConfig, error) {
	conf := AppConfig{}
	if err := envconfig.InitWithPrefix(&conf, "ABF"); err != nil {
		return conf, err
	}
	return conf, nil
}
