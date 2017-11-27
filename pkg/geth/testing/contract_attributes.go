package testing

import (
	"github.com/8thlight/vulcanizedb/pkg/core"
)

func FindAttribute(contractAttributes []core.ContractAttribute, attributeName string) *core.ContractAttribute {
	for _, contractAttribute := range contractAttributes {
		if contractAttribute.Name == attributeName {
			return &contractAttribute
		}
	}
	return nil
}
