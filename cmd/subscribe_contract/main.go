package main

import (
	"flag"

	"github.com/8thlight/vulcanizedb/cmd"
	"github.com/8thlight/vulcanizedb/pkg/core"
	"github.com/8thlight/vulcanizedb/pkg/repositories"
)

func main() {
	environment := flag.String("environment", "", "Environment name")
	contractHash := flag.String("contract-hash", "", "contract-hash=x1234")
	flag.Parse()
	config := cmd.LoadConfig(*environment)
	repository := repositories.NewPostgres(config.Database)
	repository.CreateContract(core.Contract{Hash: *contractHash})
}
