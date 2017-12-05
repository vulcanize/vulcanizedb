package testing

import (
	"path/filepath"

	"github.com/8thlight/vulcanizedb/pkg/config"
	"github.com/8thlight/vulcanizedb/pkg/core"
	"github.com/8thlight/vulcanizedb/pkg/geth"
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
	abiFilepath := filepath.Join(config.ProjectRoot(), "pkg", "geth", "testing", "sample_abi.json")
	abiFileContents, _ := geth.ReadAbiFile(abiFilepath)
	return abiFileContents
}
