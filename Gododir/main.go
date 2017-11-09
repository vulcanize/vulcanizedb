package main

import (
	"log"

	"fmt"

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

func loadConfig(environment string) config.Config {
	cfg, err := config.NewConfig(environment)
	if err != nil {
		log.Fatalf("Error loading config\n%v", err)
	}
	return *cfg
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

	p.Task("migrate", nil, func(context *do.Context) {
		environment := parseEnvironment(context)
		cfg := loadConfig(environment)
		connectString := config.DbConnectionString(cfg.Database)
		migrate := fmt.Sprintf("migrate -database '%s' -path ./db/migrations up", connectString)
		dumpSchema := fmt.Sprintf("pg_dump -O -s %s > ./db/schema.sql", cfg.Database.Name)
		context.Bash(migrate)
		context.Bash(dumpSchema)
	})

	p.Task("rollback", nil, func(context *do.Context) {
		environment := parseEnvironment(context)
		cfg := loadConfig(environment)
		connectString := config.DbConnectionString(cfg.Database)
		migrate := fmt.Sprintf("migrate -database '%s' -path ./db/migrations down 1", connectString)
		dumpSchema := fmt.Sprintf("pg_dump -O -s %s > ./db/schema.sql", cfg.Database.Name)
		context.Bash(migrate)
		context.Bash(dumpSchema)
	})

}

func main() {
	do.Godo(tasks)
}
