package config

import (
	"os"
)

type AppConfig struct {
	DittoHost           string
	DittoUsername       string
	DittoDevopsUsername string
	DittoDevopsPassword string
	DittoPassword       string
	DittoNamespace      string
	TemporalHost        string
}

func LoadConfig() AppConfig {
	return AppConfig{
		DittoHost:           os.Getenv("DITTO_HOSTPORT"),
		DittoUsername:       os.Getenv("DITTO_USERNAME"),
		DittoPassword:       os.Getenv("DITTO_PASSWORD"),
		DittoDevopsUsername: os.Getenv("DITTO_DEVOPS_USERNAME"),
		DittoDevopsPassword: os.Getenv("DITTO_DEVOPS_PASSWORD"),
		DittoNamespace:      os.Getenv("DITTO_NAMESPACE"),
		TemporalHost:        os.Getenv("TEMPORAL_HOSTPORT"),
	}
}
