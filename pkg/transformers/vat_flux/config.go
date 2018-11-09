package vat_flux

import (
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared/constants"
)

var VatFluxConfig = shared.TransformerConfig{
	TransformerName:     constants.VatFluxLabel,
	ContractAddresses:   []string{constants.VatContractAddress},
	ContractAbi:         constants.VatABI,
	Topic:               constants.VatFluxSignature,
	StartingBlockNumber: 0,
	EndingBlockNumber:   -1,
}
