package cmd

import (
	"log"

	"github.com/8thlight/vulcanizedb/pkg/config"
	"github.com/8thlight/vulcanizedb/pkg/repositories"
)

func LoadConfig(environment string) config.Config {
	cfg, err := config.NewConfig(environment)
	if err != nil {
		log.Fatalf("Error loading config\n%v", err)
	}
	return *cfg
}

func LoadPostgres(database config.Database) repositories.Postgres {
	repository, err := repositories.NewPostgres(database)
	if err != nil {
		log.Fatalf("Error loading postgres\n%v", err)
	}
	return repository
}
