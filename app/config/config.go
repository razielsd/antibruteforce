package config

import (
	"github.com/vrischmann/envconfig"
)

type AppConfig struct {
	Port          int      `envconfig:"default=8080"`
	RateLogin     int      `envconfig:"default=10"`
	RatePwd       int      `envconfig:"default=100"`
	RateIP        int      `envconfig:"default=10000"`
	Whitelist     []string `envconfig:"optional"`
	WhitelistPath string   `envconfig:"optional"`
	Blacklist     []string `envconfig:"optional"`
	BlacklistPath string   `envconfig:"optional"`
}

func GetConfig() (AppConfig, error) {
	conf := AppConfig{}
	if err := envconfig.InitWithPrefix(&conf, "ABF"); err != nil {
		return conf, err
	}
	return conf, nil
}
