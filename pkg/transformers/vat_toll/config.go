package vat_toll

import (
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared/constants"
)

func GetVatTollConfig() shared.TransformerConfig {
	return shared.TransformerConfig{
		TransformerName:     constants.VatTollLabel,
		ContractAddresses:   []string{constants.VatContractAddress()},
		ContractAbi:         constants.VatABI(),
		Topic:               constants.GetVatTollSignature(),
		StartingBlockNumber: constants.VatDeploymentBlock(),
		EndingBlockNumber:   -1,
	}
}
