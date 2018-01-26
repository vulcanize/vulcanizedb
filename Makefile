HOST_NAME =
PORT =
NAME =
CONNECT_STRING=postgresql://$(HOST_NAME):$(PORT)/$(NAME)?sslmode=disable

$(MATTESMIGRATE):
	go get -u -d github.com/mattes/migrate/cli github.com/lib/pq
	go build -tags 'postgres' -o /usr/local/bin/migrate github.com/mattes/migrate/cli

$(DEP):
	go get -u github.com/golang/dep/cmd/dep

$(GINKGO):
	go get -u github.com/onsi/ginkgo/ginkgo  

checkdbvars:
	test -n "$(HOST_NAME)" # $$HOST_NAME 
	test -n "$(PORT)" # $$PORT
	test -n "$(NAME)" # $$NAME

rollback: checkdbvars
	migrate -database $(CONNECT_STRING) -path ./db/migrations down 1
	pg_dump -O -s $(CONNECT_STRING) > db/schema.sql

migrate: $(MATTESMIGRATE) checkdbvars
	migrate -database $(CONNECT_STRING) -path ./db/migrations up
	pg_dump -O -s $(CONNECT_STRING) > db/schema.sql

import:
	test -n "$(NAME)" # $$NAME
	psql $(NAME) < db/schema.sql

dep: $(DEP)
	dep ensure

build: dep
	go build

test: $(GINKGO)
	ginkgo -r

createprivate:
	#!/bin/bash
	echo "Deleting test blockchain"
	rm -rf test_data_dir
	echo "Creating test blockchain with a new account"
	mkdir test_data_dir
	geth --dev --datadir test_data_dir --password .private_blockchain_password account new

startprivate: createprivate
	geth --datadir test_data_dir --dev --nodiscover --mine --minerthreads 1 --maxpeers 0 --verbosity 3 --unlock 0 --password .private_blockchain_password --rpc
