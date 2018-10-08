package vat_grab

import "github.com/vulcanize/vulcanizedb/pkg/transformers/shared"

var VatGrabConfig = shared.TransformerConfig{
	ContractAddresses:   []string{shared.VatContractAddress},
	ContractAbi:         shared.VatABI,
	Topics:              []string{shared.VatGrabSignature},
	StartingBlockNumber: 0,
	EndingBlockNumber:   10000000,
}
