package vat_grab

import (
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared/constants"
)

var VatGrabConfig = shared.TransformerConfig{
	TransformerName:     constants.VatGrabLabel,
	ContractAddresses:   []string{constants.VatContractAddress},
	ContractAbi:         constants.VatABI,
	Topic:               constants.VatGrabSignature,
	StartingBlockNumber: 0,
	EndingBlockNumber:   -1,
}
