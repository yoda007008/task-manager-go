package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Database struct {
		Url string `mapstructure:"url"`
	} `mapstructure:"database"`
	Server struct {
		Port string `mapstructure:"port"`
	} `mapstructure:"server"`
}

func LoadConfig(path string) (Config, error) {
	var config Config

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return config, fmt.Errorf("read config file %w", err)
	}

	expanded := os.ExpandEnv(string(data))

	viper.SetConfigType("yaml")

	if err = viper.ReadConfig(strings.NewReader(expanded)); err != nil {
		return config, fmt.Errorf("viper reading config %w", err)
	}

	if err = viper.Unmarshal(&config); err != nil {
		return config, fmt.Errorf("unmarshal config %w", err)
	}

	return config, nil
}
