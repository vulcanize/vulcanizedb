package vat_flux

import "github.com/vulcanize/vulcanizedb/pkg/transformers/shared"

var VatFluxConfig = shared.TransformerConfig{
	ContractAddresses:   []string{shared.VatContractAddress},
	ContractAbi:         shared.VatABI,
	Topics:              []string{shared.VatFluxSignature},
	StartingBlockNumber: 0,
	EndingBlockNumber:   10000000,
}
