package graphql_server_test

import (
	"log"

	"encoding/json"

	"context"

	"github.com/neelance/graphql-go"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/config"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/filters"
	"github.com/vulcanize/vulcanizedb/pkg/graphql_server"
	"github.com/vulcanize/vulcanizedb/pkg/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/repositories/postgres"
)

func formatJSON(data []byte) []byte {
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		log.Fatalf("invalid JSON: %s", err)
	}
	formatted, err := json.Marshal(v)
	if err != nil {
		log.Fatal(err)
	}
	return formatted
}

var _ = Describe("GraphQL", func() {
	var cfg config.Config
	var repository repositories.Repository

	BeforeEach(func() {

		cfg, _ = config.NewConfig("private")
		node := core.Node{GenesisBlock: "GENESIS", NetworkId: 1, Id: "x123", ClientName: "geth"}
		repository = postgres.BuildRepository(node)
		e := repository.AddFilter(filters.LogFilter{
			Name:      "TestFilter1",
			FromBlock: 1,
			ToBlock:   10,
			Address:   "0x123456789",
			Topics:    core.Topics{0: "event=1", 2: "event=2"},
		})
		if e != nil {
			log.Fatal(e)
		}
		f, e := repository.GetFilter("TestFilter1")
		if e != nil {
			log.Println(f)
			log.Fatal(e)
		}
	})

	It("Queries example schema for specific log filter", func() {
		var variables map[string]interface{}
		r := graphql_server.NewResolver(repository)
		var schema = graphql.MustParseSchema(graphql_server.Schema, r)
		response := schema.Exec(context.Background(),
			`{
                            logFilter(name: "TestFilter1") {
                                name
                                fromBlock
                                toBlock
                                address
                                topics
                             }
                           }`,
			"",
			variables)
		expected := `{
                        "logFilter": {
						    "name": "TestFilter1", 
						    "fromBlock": 1, 
						    "toBlock": 10,
                            "address": "0x123456789",
                            "topics": ["event=1", null, "event=2", null]
						}
						 }`
		var v interface{}
		if len(response.Errors) != 0 {
			log.Fatal(response.Errors)
		}
		err := json.Unmarshal(response.Data, &v)
		Expect(err).ToNot(HaveOccurred())
		a := formatJSON(response.Data)
		e := formatJSON([]byte(expected))
		Expect(a).To(Equal(e))
	})
})
