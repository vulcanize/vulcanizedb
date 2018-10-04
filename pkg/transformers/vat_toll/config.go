package vat_toll

import "github.com/vulcanize/vulcanizedb/pkg/transformers/shared"

var VatTollConfig = shared.TransformerConfig{
	ContractAddresses:   []string{shared.VatContractAddress},
	ContractAbi:         shared.VatABI,
	Topics:              []string{shared.VatTollSignature},
	StartingBlockNumber: 0,
	EndingBlockNumber:   10000000,
}
