package cmd

import (
	"log"

	"path/filepath"

	"fmt"

	"math/big"

	"github.com/vulcanize/vulcanizedb/pkg/config"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/geth"
	"github.com/vulcanize/vulcanizedb/pkg/repositories"
)

func LoadConfig(environment string) config.Config {
	cfg, err := config.NewConfig(environment)
	if err != nil {
		log.Fatalf("Error loading config\n%v", err)
	}
	return cfg
}

func LoadPostgres(database config.Database, node core.Node) repositories.Postgres {
	repository, err := repositories.NewPostgres(database, node)
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

func GetAbi(abiFilepath string, contractHash string) string {
	var contractAbiString string
	if abiFilepath != "" {
		contractAbiString = ReadAbiFile(abiFilepath)
	} else {
		etherscan := geth.NewEtherScanClient("https://api.etherscan.io")
		fmt.Println("No ABI supplied. Retrieving ABI from Etherscan")
		contractAbiString, _ = etherscan.GetAbi(contractHash)
	}
	_, err := geth.ParseAbi(contractAbiString)
	if err != nil {
		log.Fatalln("Invalid ABI")
	}
	return contractAbiString
}

func RequestedBlockNumber(blockNumber *int64) *big.Int {
	var _blockNumber *big.Int
	if *blockNumber == -1 {
		_blockNumber = nil
	} else {
		_blockNumber = big.NewInt(*blockNumber)
	}
	return _blockNumber
}
