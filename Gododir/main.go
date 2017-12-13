package main

import (
	"log"

	"fmt"

	"github.com/8thlight/vulcanizedb/cmd"
	"github.com/8thlight/vulcanizedb/pkg/config"
	do "gopkg.in/godo.v2"
)

func parseEnvironment(context *do.Context) string {
	environment := context.Args.MayString("", "environment", "env", "e")
	if environment == "" {
		log.Fatalln("--environment required")
	}
	return environment
}

func tasks(p *do.Project) {

	p.Task("run", nil, func(context *do.Context) {
		environment := parseEnvironment(context)
		context.Start(`go run main.go --environment={{.environment}}`,
			do.M{"environment": environment, "$in": "cmd/run"})
	})

	p.Task("populateBlocks", nil, func(context *do.Context) {
		environment := parseEnvironment(context)
		startingNumber := context.Args.MayInt(-1, "starting-number")
		if startingNumber < 0 {
			log.Fatalln("--starting-number required")
		}
		context.Start(`go run main.go --environment={{.environment}} --starting-number={{.startingNumber}}`,
			do.M{"environment": environment, "startingNumber": startingNumber, "$in": "cmd/populate_blocks"})
	})

	p.Task("getLogs", nil, func(context *do.Context) {
		environment := parseEnvironment(context)
		blockNumber := context.Args.MayInt(-1, "block-number", "b")
		contractHash := context.Args.MayString("", "contract-hash", "c")
		if contractHash == "" {
			log.Fatalln("--contract-hash required")
		}
		context.Start(`go run main.go --environment={{.environment}} --contract-hash={{.contractHash}} --block-number={{.blockNumber}}`,
			do.M{
				"environment":  environment,
				"contractHash": contractHash,
				"blockNumber":  blockNumber,
				"$in":          "cmd/get_logs",
			})
	})

	p.Task("watchContract", nil, func(context *do.Context) {
		environment := parseEnvironment(context)
		contractHash := context.Args.MayString("", "contract-hash", "c")
		abiFilepath := context.Args.MayString("", "abi-filepath", "a")
		if contractHash == "" {
			log.Fatalln("--contract-hash required")
		}
		context.Start(`go run main.go --environment={{.environment}} --contract-hash={{.contractHash}} --abi-filepath={{.abiFilepath}}`,
			do.M{
				"environment":  environment,
				"contractHash": contractHash,
				"abiFilepath":  abiFilepath,
				"$in":          "cmd/watch_contract",
			})
	})

	p.Task("migrate", nil, func(context *do.Context) {
		environment := parseEnvironment(context)
		cfg := cmd.LoadConfig(environment)
		connectString := config.DbConnectionString(cfg.Database)
		migrate := fmt.Sprintf("migrate -database '%s' -path ./db/migrations up", connectString)
		dumpSchema := fmt.Sprintf("pg_dump -O -s %s > ./db/schema.sql", cfg.Database.Name)
		context.Bash(migrate)
		context.Bash(dumpSchema)
	})

	p.Task("rollback", nil, func(context *do.Context) {
		environment := parseEnvironment(context)
		cfg := cmd.LoadConfig(environment)
		connectString := config.DbConnectionString(cfg.Database)
		migrate := fmt.Sprintf("migrate -database '%s' -path ./db/migrations down 1", connectString)
		dumpSchema := fmt.Sprintf("pg_dump -O -s %s > ./db/schema.sql", cfg.Database.Name)
		context.Bash(migrate)
		context.Bash(dumpSchema)
	})

	p.Task("showContractSummary", nil, func(context *do.Context) {
		environment := parseEnvironment(context)
		contractHash := context.Args.MayString("", "contract-hash", "c")
		blockNumber := context.Args.MayInt(-1, "block-number", "b")
		if contractHash == "" {
			log.Fatalln("--contract-hash required")
		}
		context.Start(`go run main.go --environment={{.environment}} --contract-hash={{.contractHash}} --block-number={{.blockNumber}}`,
			do.M{"environment": environment,
				"contractHash": contractHash,
				"blockNumber":  blockNumber,
				"$in":          "cmd/show_contract_summary"})
	})

}

func main() {
	do.Godo(tasks)
}
