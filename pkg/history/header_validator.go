package history

import (
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore"
	"github.com/vulcanize/vulcanizedb/pkg/transformers"
)

type HeaderValidator struct {
	blockChain       core.BlockChain
	headerRepository datastore.HeaderRepository
	windowSize       int
}

func NewHeaderValidator(blockChain core.BlockChain, repository datastore.HeaderRepository, windowSize int) HeaderValidator {
	return HeaderValidator{
		blockChain:       blockChain,
		headerRepository: repository,
		windowSize:       windowSize,
	}
}

func (validator HeaderValidator) ValidateHeaders() ValidationWindow {
	window := MakeValidationWindow(validator.blockChain, validator.windowSize)
	blockNumbers := MakeRange(window.LowerBound, window.UpperBound)
	RetrieveAndUpdateHeaders(validator.blockChain, validator.headerRepository, blockNumbers, []transformers.Transformer{})
	return window
}
