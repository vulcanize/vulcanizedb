package vat_flux

import (
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared/constants"
)

func GetVatFluxConfig() shared.TransformerConfig {
	return shared.TransformerConfig{
		TransformerName:     constants.VatFluxLabel,
		ContractAddresses:   []string{constants.VatContractAddress()},
		ContractAbi:         constants.VatABI(),
		Topic:               constants.GetVatFluxSignature(),
		StartingBlockNumber: constants.VatDeploymentBlock(),
		EndingBlockNumber:   -1,
	}
}
