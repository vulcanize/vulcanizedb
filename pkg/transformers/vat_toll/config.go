package vat_toll

import (
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared/constants"
)

var VatTollConfig = shared.TransformerConfig{
	TransformerName:     constants.VatTollLabel,
	ContractAddresses:   []string{constants.VatContractAddress},
	ContractAbi:         constants.VatABI,
	Topic:               constants.VatTollSignature,
	StartingBlockNumber: 0,
	EndingBlockNumber:   10000000,
}
