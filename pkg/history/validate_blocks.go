package history

import (
	"io"
	"text/template"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore"
)

const WindowTemplate = `Validating Blocks
|{{.LowerBound}}|-- Validation Window --|{{.UpperBound}}| ({{.UpperBound}}:HEAD)

`

var ParsedWindowTemplate = *template.Must(template.New("window").Parse(WindowTemplate))

type BlockValidator struct {
	blockchain            core.Blockchain
	blockRepository       datastore.BlockRepository
	windowSize            int
	parsedLoggingTemplate template.Template
}

func NewBlockValidator(blockchain core.Blockchain, blockRepository datastore.BlockRepository, windowSize int) *BlockValidator {
	return &BlockValidator{
		blockchain,
		blockRepository,
		windowSize,
		ParsedWindowTemplate,
	}
}

func (bv BlockValidator) ValidateBlocks() ValidationWindow {
	window := MakeValidationWindow(bv.blockchain, bv.windowSize)
	blockNumbers := MakeRange(window.LowerBound, window.UpperBound)
	RetrieveAndUpdateBlocks(bv.blockchain, bv.blockRepository, blockNumbers)
	lastBlock := bv.blockchain.LastBlock().Int64()
	bv.blockRepository.SetBlocksStatus(lastBlock)
	return window
}

func (bv BlockValidator) Log(out io.Writer, window ValidationWindow) {
	bv.parsedLoggingTemplate.Execute(out, window)
}

type ValidationWindow struct {
	LowerBound int64
	UpperBound int64
}

func (window ValidationWindow) Size() int {
	return int(window.UpperBound - window.LowerBound)
}

func MakeValidationWindow(blockchain core.Blockchain, windowSize int) ValidationWindow {
	upperBound := blockchain.LastBlock().Int64()
	lowerBound := upperBound - int64(windowSize)
	return ValidationWindow{lowerBound, upperBound}
}

func MakeRange(min, max int64) []int64 {
	a := make([]int64, max-min)
	for i := range a {
		a[i] = min + int64(i)
	}
	return a
}
