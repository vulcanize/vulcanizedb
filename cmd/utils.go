package cmd

import (
	"log"

	"path/filepath"

	"github.com/8thlight/vulcanizedb/pkg/config"
	"github.com/8thlight/vulcanizedb/pkg/geth"
	"github.com/8thlight/vulcanizedb/pkg/repositories"
)

func LoadConfig(environment string) config.Config {
	cfg, err := config.NewConfig(environment)
	if err != nil {
		log.Fatalf("Error loading config\n%v", err)
	}
	return cfg
}

func LoadPostgres(database config.Database) repositories.Postgres {
	repository, err := repositories.NewPostgres(database)
	if err != nil {
		log.Fatalf("Error loading postgres\n%v", err)
	}
	return repository
}

func ReadAbiFile(abiFilepath string) string {
	if !filepath.IsAbs(abiFilepath) {
		abiFilepath = filepath.Join(config.ProjectRoot(), abiFilepath)
	}
	abi, err := geth.ReadAbiFile(abiFilepath)
	if err != nil {
		log.Fatalf("Error reading ABI file at \"%s\"\n %v", abiFilepath, err)
	}
	return abi
}
