package vat_grab

import (
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared/constants"
)

func GetVatGrabConfig() shared.TransformerConfig {
	return shared.TransformerConfig{
		TransformerName:     constants.VatGrabLabel,
		ContractAddresses:   []string{constants.VatContractAddress()},
		ContractAbi:         constants.VatABI(),
		Topic:               constants.GetVatGrabSignature(),
		StartingBlockNumber: constants.VatDeploymentBlock(),
		EndingBlockNumber:   -1,
	}
}
