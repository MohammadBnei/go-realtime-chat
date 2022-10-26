package config

import (
	"flag"
	"fmt"
	"log"

	"github.com/spf13/viper"
)

func ParseFlags() string {
	ConfigFilePath := flag.String("config", ".", "path to config file")

	flag.Parse()
	return *ConfigFilePath
}

// Flags returned from the function
var ConfigFilePath = ParseFlags()

func ParseConfig() config {
	readConfig := config{}

	log.Println("parsing config file")
	viper.SetConfigType("yml")
	viper.AddConfigPath(".")    // optionally look for config in the working directory
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	viper.AutomaticEnv()

	err = viper.Unmarshal(&readConfig)
	if err != nil {
		log.Fatal("error unmarshing config file in struct : ", err)
	}

	readConfig.ServerConfig.Port = viper.Get("port").(string)

	log.Println("Env variable parsed successfully")

	return readConfig
}

var parsedConfig = ParseConfig()
