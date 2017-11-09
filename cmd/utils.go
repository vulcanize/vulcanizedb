package cmd

import (
	"log"

	"github.com/8thlight/vulcanizedb/pkg/config"
)

func LoadConfig(environment string) config.Config {
	cfg, err := config.NewConfig(environment)
	if err != nil {
		log.Fatalf("Error loading config\n%v", err)
	}
	return *cfg
}
