package vat_flux

import "github.com/vulcanize/vulcanizedb/pkg/transformers/shared"

var VatFluxConfig = shared.SingleTransformerConfig{
	TransformerName:     shared.VatFluxLabel,
	ContractAddresses:   []string{shared.VatContractAddress},
	ContractAbi:         shared.VatABI,
	Topic:               shared.VatFluxSignature,
	StartingBlockNumber: 0,
	EndingBlockNumber:   10000000,
}
