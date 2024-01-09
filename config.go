package main

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	MySQL MySQLConfig `mapstructure:"mysql"`
	Redis RedisConfig `mapstructure:"redis"`
	Port  string      `mapstructure:"port"`
}

type MySQLConfig struct {
	ConnectionString string `mapstructure:"connection_string"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Password string `mapstructure:"password"`
}

var AppConfig *Config

func LoadAppConfig() {
	log.Println("Loading Server Configurations...")
	viper.AddConfigPath(".")

	env := os.Getenv("env")
	if env == "docker" {
		viper.SetConfigName("config-docker")
	} else {
		viper.SetConfigName("config")
	}

	viper.SetConfigType("json")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
	err = viper.Unmarshal(&AppConfig)
	if err != nil {
		log.Fatal(err)
	}
}
