package config

import (
	cfg "github.com/lee0720/nuwa/pkg/config"
)

type ConfigType struct {
	SecondaryMarketConfig cfg.MySQLConfiguration
	InfluxConfig          cfg.InfluxConfiguration
}

var CONFIG ConfigType
