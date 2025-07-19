package config

import (
	"os"
)

type AppConfig struct {
	DittoHost     string
	DittoUsername string
	DittoPassword string
	TemporalHost  string
}

func LoadConfig() AppConfig {
	return AppConfig{
		DittoHost:     os.Getenv("DITTO_HOSTPORT"),
		DittoUsername: os.Getenv("DITTO_USERNAME"),
		DittoPassword: os.Getenv("DITTO_PASSWORD"),
		TemporalHost:  os.Getenv("TEMPORAL_HOSTPORT"),
	}
}
