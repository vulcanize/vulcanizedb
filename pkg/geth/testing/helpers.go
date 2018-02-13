package testing

import (
	"log"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/geth"
	"github.com/vulcanize/vulcanizedb/test_config"
)

func FindAttribute(contractAttributes core.ContractAttributes, attributeName string) *core.ContractAttribute {
	for _, contractAttribute := range contractAttributes {
		if contractAttribute.Name == attributeName {
			return &contractAttribute
		}
	}
	return nil
}

func SampleContract() core.Contract {
	return core.Contract{
		Abi:  sampleAbiFileContents(),
		Hash: "0xd26114cd6EE289AccF82350c8d8487fedB8A0C07",
	}
}

func sampleAbiFileContents() string {
	abiFileContents, err := geth.ReadAbiFile(test_config.ABIFilePath + "sample_abi.json")
	if err != nil {
		log.Fatal(err)
	}
	return abiFileContents
}
