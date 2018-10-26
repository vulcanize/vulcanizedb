package vat_grab

import "github.com/vulcanize/vulcanizedb/pkg/transformers/shared"

var VatGrabConfig = shared.SingleTransformerConfig{
	TransformerName:     shared.VatGrabLabel,
	ContractAddresses:   []string{shared.VatContractAddress},
	ContractAbi:         shared.VatABI,
	Topic:               shared.VatGrabSignature,
	StartingBlockNumber: 0,
	EndingBlockNumber:   10000000,
}
