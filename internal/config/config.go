package config

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Port     int                  `mapstructure:"port"`
	Products []map[string]Product `mapstructure:"products"`
}

type Product struct {
	Host    string `mapstructure:"host"`
	Version string `mapstructure:"version"`
}

func SetConfig(configPath string) Config {
	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Printf("Config file not found at %s", configPath)
		log.Printf("Please provide a valid config file using the --config flag")
		log.Printf("Example config file format:")
		log.Printf(`
port: 2112
products:
  - ubuntu:
      host: localhost
      version: "22.04"
  - nodejs:
      host: localhost
      version: "16"`)
		os.Exit(1)
	}

	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Error reading config file: %v", err)
		panic(fmt.Sprintf("Error reading config file: %v", err))
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Printf("Error unmarshalling config: %v", err)
		panic(fmt.Sprintf("Error reading config file: %v", err))
	}

	if config.Port == 0 {
		log.Printf("Port must be specified in config file")
		panic("Error reading config file: port must be specified")
	}

	if len(config.Products) == 0 {
		log.Printf("At least one product must be specified in config file")
		panic("Error reading config file: at least one product must be specified")
	}

	return config
}
