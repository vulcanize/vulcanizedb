package pit_file

import (
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
)

var PitFileConfig = shared.TransformerConfig{
	ContractAddress:     "0xff3f2400f1600f3f493a9a92704a29b96795af1a", // temporary address from Ganache deploy
	ContractAbi:         shared.PitABI,
	Topics:              []string{shared.PitFileSignatureOne},
	StartingBlockNumber: 0,
	EndingBlockNumber:   100,
}
