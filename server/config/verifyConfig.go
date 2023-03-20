package config

import "log"

func VerifyConfig() config {

	if parsedConfig.Port == "" {
		log.Fatal("No port specified")
	}

	return parsedConfig

}

var Config = VerifyConfig()
