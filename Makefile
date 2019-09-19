BIN = $(GOPATH)/bin
BASE = $(GOPATH)/src/$(PACKAGE)
PKGS = go list ./... | grep -v "^vendor/"

# Tools
## Testing library
GINKGO = $(BIN)/ginkgo
$(BIN)/ginkgo:
	go get -u github.com/onsi/ginkgo/ginkgo

## Migration tool
GOOSE = $(BIN)/goose
$(BIN)/goose:
	go get -u -d github.com/pressly/goose/cmd/goose
	go build -tags='no_mysql no_sqlite' -o $(BIN)/goose github.com/pressly/goose/cmd/goose

## Source linter
LINT = $(BIN)/golint
$(BIN)/golint:
	go get -u golang.org/x/lint/golint

## Combination linter
METALINT = $(BIN)/gometalinter.v2
$(BIN)/gometalinter.v2:
	go get -u gopkg.in/alecthomas/gometalinter.v2
	$(METALINT) --install


.PHONY: installtools
installtools: | $(LINT) $(GOOSE) $(GINKGO)
	echo "Installing tools"

.PHONY: metalint
metalint: | $(METALINT)
	$(METALINT) ./... --vendor \
	--fast \
	--exclude="exported (function)|(var)|(method)|(type).*should have comment or be unexported" \
	--format="{{.Path.Abs}}:{{.Line}}:{{if .Col}}{{.Col}}{{end}}:{{.Severity}}: {{.Message}} ({{.Linter}})"

.PHONY: lint
lint:
	$(LINT) $$($(PKGS)) | grep -v -E "exported (function)|(var)|(method)|(type).*should have comment or be unexported"

#Test
TEST_DB = vulcanize_testing
TEST_CONNECT_STRING = postgresql://localhost:5432/$(TEST_DB)?sslmode=disable

.PHONY: test
test: | $(GINKGO) $(LINT)
	go vet ./...
	go fmt ./...
	dropdb --if-exists $(TEST_DB)
	createdb $(TEST_DB)
	cd db/migrations;\
    	$(GOOSE) postgres "$(TEST_CONNECT_STRING)" up
	cd db/migrations/;\
		$(GOOSE) postgres "$(TEST_CONNECT_STRING)" reset
	make migrate NAME=$(TEST_DB)
	$(GINKGO) -r --skipPackage=integration_tests,integration

.PHONY: integrationtest
integrationtest: | $(GINKGO) $(LINT)
	go vet ./...
	go fmt ./...
	dropdb --if-exists $(TEST_DB)
	createdb $(TEST_DB)
	cd db/migrations;\
	    $(GOOSE) postgres "$(TEST_CONNECT_STRING)" up
	cd db/migrations/;\
    	$(GOOSE) postgres "$(TEST_CONNECT_STRING)" reset
	make migrate NAME=$(TEST_DB)
	$(GINKGO) -r integration_test/

build:
	go fmt ./...
	go build

#Database
HOST_NAME = localhost
PORT = 5432
NAME =
USER = postgres
CONNECT_STRING=postgresql://$(USER)@$(HOST_NAME):$(PORT)/$(NAME)?sslmode=disable

# Parameter checks
## Check that DB variables are provided
.PHONY: checkdbvars
checkdbvars:
	test -n "$(HOST_NAME)" # $$HOST_NAME
	test -n "$(PORT)" # $$PORT
	test -n "$(NAME)" # $$NAME
	@echo $(CONNECT_STRING)

## Check that the migration variable (id/timestamp) is provided
.PHONY: checkmigration
checkmigration:
	test -n "$(MIGRATION)" # $$MIGRATION

# Check that the migration name is provided
.PHONY: checkmigname
checkmigname:
	test -n "$(NAME)" # $$NAME

# Migration operations
## Rollback the last migration
.PHONY: rollback
rollback: $(GOOSE) checkdbvars
	cd db/migrations;\
	  $(GOOSE) postgres "$(CONNECT_STRING)" down
	pg_dump -O -s $(CONNECT_STRING) > db/schema.sql


## Rollbackt to a select migration (id/timestamp)
.PHONY: rollback_to
rollback_to: $(GOOSE) checkmigration checkdbvars
	cd db/migrations;\
	  $(GOOSE) postgres "$(CONNECT_STRING)" down-to "$(MIGRATION)"

## Apply all migrations not already run
.PHONY: migrate
migrate: $(GOOSE) checkdbvars
	cd db/migrations;\
	  $(GOOSE) postgres "$(CONNECT_STRING)" up
	pg_dump -O -s $(CONNECT_STRING) > db/schema.sql

## Create a new migration file
.PHONY: new_migration
new_migration: $(GOOSE) checkmigname
	cd db/migrations;\
	  $(GOOSE) create $(NAME) sql

## Check which migrations are applied at the moment
.PHONY: migration_status
migration_status: $(GOOSE) checkdbvars
	cd db/migrations;\
	  $(GOOSE) postgres "$(CONNECT_STRING)" status

# Convert timestamped migrations to versioned (to be run in CI);
# merge timestamped files to prevent conflict
.PHONY: version_migrations
version_migrations:
	cd db/migrations; $(GOOSE) fix

# Import a psql schema to the database
.PHONY: import
import:
	test -n "$(NAME)" # $$NAME
	psql $(NAME) < db/schema.sql


# Docker actions
## Rinkeby docker environment
RINKEBY_COMPOSE_FILE=dockerfiles/rinkeby/docker-compose.yml

.PHONY: rinkeby_env_up
rinkeby_env_up:
	docker-compose -f $(RINKEBY_COMPOSE_FILE) up -d geth
	docker-compose -f $(RINKEBY_COMPOSE_FILE) up --build migrations
	docker-compose -f $(RINKEBY_COMPOSE_FILE) up -d --build vulcanizedb

.PHONY: rinkeby_env_deploy
rinkeby_env_deploy:
	docker-compose -f $(RINKEBY_COMPOSE_FILE) up -d --build vulcanizedb

.PHONY: dev_env_migrate
rinkeby_env_migrate:
	docker-compose -f $(RINKEBY_COMPOSE_FILE) up --build migrations

.PHONY: rinkeby_env_down
rinkeby_env_down:
	docker-compose -f $(RINKEBY_COMPOSE_FILE) down
