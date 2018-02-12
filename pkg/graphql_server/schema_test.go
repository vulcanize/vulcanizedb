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
	var graphQLRepositories graphql_server.GraphQLRepositories

	BeforeEach(func() {

		cfg, _ = config.NewConfig("private")
		node := core.Node{GenesisBlock: "GENESIS", NetworkId: 1, Id: "x123", ClientName: "geth"}
		db := postgres.NewTestDB(node)
		blockRepository := &postgres.BlockRepository{DB: db}
		logRepository := &postgres.LogRepository{DB: db}
		filterRepository := &postgres.FilterRepository{DB: db}
		watchedEventRepository := &postgres.WatchedEventRepository{DB: db}
		graphQLRepositories = graphql_server.GraphQLRepositories{
			WatchedEventRepository: watchedEventRepository,
			BlockRepository:        blockRepository,
			LogRepository:          logRepository,
			FilterRepository:       filterRepository,
		}

		err := graphQLRepositories.CreateFilter(filters.LogFilter{
			Name:      "TestFilter1",
			FromBlock: 1,
			ToBlock:   10,
			Address:   "0x123456789",
			Topics:    core.Topics{0: "topic=1", 2: "topic=2"},
		})
		if err != nil {
			log.Fatal(err)
		}
		filter, err := graphQLRepositories.GetFilter("TestFilter1")
		if err != nil {
			log.Println(filter)
			log.Fatal(err)
		}

		matchingEvent := core.Log{
			BlockNumber: 5,
			TxHash:      "0xTX1",
			Address:     "0x123456789",
			Topics:      core.Topics{0: "topic=1", 2: "topic=2"},
			Index:       0,
			Data:        "0xDATADATADATA",
		}
		nonMatchingEvent := core.Log{
			BlockNumber: 5,
			TxHash:      "0xTX2",
			Address:     "0xOTHERADDRESS",
			Topics:      core.Topics{0: "topic=1", 2: "topic=2"},
			Index:       0,
			Data:        "0xDATADATADATA",
		}
		err = graphQLRepositories.CreateLogs([]core.Log{matchingEvent, nonMatchingEvent})
		if err != nil {
			log.Fatal(err)
		}
	})

	It("Queries example schema for specific log filter", func() {
		var variables map[string]interface{}
		resolver := graphql_server.NewResolver(graphQLRepositories)
		var schema = graphql.MustParseSchema(graphql_server.Schema, resolver)
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
                            "topics": ["topic=1", null, "topic=2", null]
						}
						 }`
		var v interface{}
		if len(response.Errors) != 0 {
			log.Fatal(response.Errors)
		}
		err := json.Unmarshal(response.Data, &v)
		Expect(err).ToNot(HaveOccurred())
		actualJSON := formatJSON(response.Data)
		expectedJSON := formatJSON([]byte(expected))
		Expect(actualJSON).To(Equal(expectedJSON))
	})

	It("Queries example schema for specific watched event log", func() {
		var variables map[string]interface{}

		resolver := graphql_server.NewResolver(graphQLRepositories)
		var schema = graphql.MustParseSchema(graphql_server.Schema, resolver)
		response := schema.Exec(context.Background(),
			`{
                           watchedEvents(name: "TestFilter1") {
                            total
                            watchedEvents{
                                name
                                blockNumber
                                address
                                tx_hash
                                topic0
                                topic1
                                topic2
                                topic3
                                data
                              }
                            }
                        }`,
			"",
			variables)
		expected := `{
	                  "watchedEvents":
	                     {
                            "total": 1,
                            "watchedEvents": [
                                {"name":"TestFilter1",
                                 "blockNumber": 5,
                                 "address": "0x123456789",
                                 "tx_hash": "0xTX1",
                                 "topic0": "topic=1",
                                 "topic1": "",
                                 "topic2": "topic=2",
                                 "topic3": "",
                                 "data": "0xDATADATADATA"
                                }
                            ]
	                     }
	               }`
		var v interface{}
		if len(response.Errors) != 0 {
			log.Fatal(response.Errors)
		}
		err := json.Unmarshal(response.Data, &v)
		Expect(err).ToNot(HaveOccurred())
		actualJSON := formatJSON(response.Data)
		expectedJSON := formatJSON([]byte(expected))
		Expect(actualJSON).To(Equal(expectedJSON))
	})
})
