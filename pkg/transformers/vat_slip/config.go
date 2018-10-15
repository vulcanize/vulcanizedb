package vat_slip

import "github.com/vulcanize/vulcanizedb/pkg/transformers/shared"

var VatSlipConfig = shared.TransformerConfig{
	ContractAddresses:   []string{shared.VatContractAddress},
	ContractAbi:         shared.VatABI,
	Topics:              []string{shared.VatSlipSignature},
	StartingBlockNumber: 0,
	EndingBlockNumber:   10000000,
}
