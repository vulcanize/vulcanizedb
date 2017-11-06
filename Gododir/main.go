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

func tasks(p *do.Project) {

	p.Task("run", nil, func(context *do.Context) {
		environment := parseEnvironment(context)
		context.Start(`go run main.go --environment={{.environment}}`,
			do.M{"environment": environment, "$in": "cmd/run"})
	})

	p.Task("migrate", nil, func(context *do.Context) {
		environment := parseEnvironment(context)
		cfg := config.NewConfig(environment)
		connectString := config.DbConnectionString(cfg.Database)
		migrate := fmt.Sprintf("migrate -database '%s' -path ./db/migrations up", connectString)
		dumpSchema := fmt.Sprintf("pg_dump -O -s %s > ./db/schema.sql", cfg.Database.Name)
		context.Bash(migrate)
		context.Bash(dumpSchema)
	})

	p.Task("rollback", nil, func(context *do.Context) {
		environment := parseEnvironment(context)
		cfg := config.NewConfig(environment)
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
