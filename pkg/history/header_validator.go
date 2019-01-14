package history

import (
	log "github.com/sirupsen/logrus"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore"
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

func (validator HeaderValidator) ValidateHeaders() (ValidationWindow, error) {
	window := MakeValidationWindow(validator.blockChain, validator.windowSize)
	blockNumbers := MakeRange(window.LowerBound, window.UpperBound)
	_, err := RetrieveAndUpdateHeaders(validator.blockChain, validator.headerRepository, blockNumbers)
	if err != nil {
		log.Error("Error in ValidateHeaders: ", err)
		return ValidationWindow{}, err
	}
	return window, nil
}
