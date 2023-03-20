package config

import (
	"flag"
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

	viper.SetConfigType("yml") // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(".")   // optionally look for config in the working directory
	viper.BindEnv("Port")
	err := viper.ReadInConfig()

	if err != nil {
		log.Fatal("error parsing config file : ", err)
	}

	viper.AutomaticEnv()

	err = viper.Unmarshal(&readConfig)
	if err != nil {
		log.Fatal("error unmarshing config file in struct : ", err)
	}

	log.Println("Env variable parsed successfully")

	return readConfig
}

var parsedConfig = ParseConfig()
