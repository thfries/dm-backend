package config

import (
	"os"
)

type AppConfig struct {
	DittoHost    string
	TemporalHost string
}

func LoadConfig() AppConfig {
	return AppConfig{
		DittoHost:    os.Getenv("DITTO_HOSTPORT"),
		TemporalHost: os.Getenv("TEMPORAL_HOSTPORT"),
	}
}
