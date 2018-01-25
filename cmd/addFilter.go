package cmd

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/8thlight/vulcanizedb/pkg/filters"
	"github.com/8thlight/vulcanizedb/pkg/geth"
	"github.com/8thlight/vulcanizedb/utils"
	"github.com/spf13/cobra"
)

// addFilterCmd represents the addFilter command
var addFilterCmd = &cobra.Command{
	Use:   "addFilter",
	Short: "Adds event filter to vulcanizedb",
	Long: `An event filter is added to the vulcanize_db. 
All events matching the filter conitions will be tracked 
in vulcanizedb. 

vulcanizedb addFilter --config config.toml --filter-filepath filter.json

The event filters are expected to match
the format described in the ethereum RPC wiki:

https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_newfilter

[{
  "fromBlock": "0x1",
  "toBlock": "0x2",
  "address": "0x8888f1f195afa192cfee860698584c030f4c9db1",
  "topics": ["0x000000000000000000000000a94f5374fce5edbc8e2a8697c15331677e6ebf0b", 
             null, 
             "0x000000000000000000000000a94f5374fce5edbc8e2a8697c15331677e6ebf0b", 
             "0x0000000000000000000000000aff3454fce5edbc8cca8697c15331677e6ebccc"]
}]
`,
	Run: func(cmd *cobra.Command, args []string) {
		addFilter()
	},
}

var filterFilepath string

func init() {
	rootCmd.AddCommand(addFilterCmd)

	addFilterCmd.PersistentFlags().StringVar(&filterFilepath, "filter-filepath", "", "path/to/filter.json")
	addFilterCmd.MarkFlagRequired("filter-filepath")
}

func addFilter() {
	if filterFilepath == "" {
		log.Fatal("filter-filepath required")
	}
	var logFilters filters.LogFilters
	blockchain := geth.NewBlockchain(ipc)
	repository := utils.LoadPostgres(databaseConfig, blockchain.Node())
	absFilePath := utils.AbsFilePath(filterFilepath)
	logFilterBytes, err := ioutil.ReadFile(absFilePath)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(logFilterBytes, &logFilters)
	if err != nil {
		log.Fatal(err)
	}
	for _, filter := range logFilters {
		err = repository.AddFilter(filter)
		if err != nil {
			log.Fatal(err)
		}
	}
}
