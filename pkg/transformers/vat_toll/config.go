package vat_toll

import "github.com/vulcanize/vulcanizedb/pkg/transformers/shared"

var VatTollConfig = shared.SingleTransformerConfig{
	TransformerName:     shared.VatTollLabel,
	ContractAddresses:   []string{shared.VatContractAddress},
	ContractAbi:         shared.VatABI,
	Topic:               shared.VatTollSignature,
	StartingBlockNumber: 0,
	EndingBlockNumber:   10000000,
}
