package config

import (
	"flag"
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Env         string
	Application ApplicationConfig `json:"Application"`
	Store       StoreConfig       `json:"Store"`
}

type ApplicationConfig struct {
	HostUrl     string        `json:"hostUrl"`
	Secret      string        `json:"secret"`
	TokenExpiry time.Duration `json:"tokenExpiry"`
}

type StoreConfig struct {
	Connection     string `json:"connection"`
	MigrationsPath string `json:"migrationsPath"`
}

func Load() (*Config, error) {
	var env string
	flag.StringVar(&env, "env", "development", "environment")
	flag.Parse()

	configName := fmt.Sprintf("config.%s", env)
	var v = viper.New()
	v.SetConfigName(configName)
	v.SetConfigType("json")
	v.AddConfigPath(".")
	err := v.ReadInConfig()
	if err != nil {
		return nil, err
	}
	var conf Config
	if err = v.Unmarshal(&conf); err != nil {
		return nil, err
	}
	conf.Env = env
	return &conf, nil
}
