package config

import (
	"flag"
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Env         string
	Application applicationConfig `json:"Application"`
	Store       storeConfig       `json:"Store"`
}

type applicationConfig struct {
	HostUrl string `json:"hostUrl"`
	Secret  string `json:"secret"`
}

type storeConfig struct {
	ConnString     string `json:"connectionString"`
	MigrationsPath string `json:"migrationsPath"`
}

type xConfig struct {
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
