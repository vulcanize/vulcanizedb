package vat_tune

import "github.com/vulcanize/vulcanizedb/pkg/transformers/shared"

var VatTuneConfig = shared.TransformerConfig{
	ContractAddresses:   []string{shared.VatContractAddress},
	ContractAbi:         shared.VatABI,
	Topics:              []string{shared.VatTuneSignature},
	StartingBlockNumber: 0,
	EndingBlockNumber:   10000000,
}
