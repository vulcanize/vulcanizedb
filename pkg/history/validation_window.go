package history

import (
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"io"
	"text/template"
)

const WindowTemplate = `Validating Blocks
|{{.LowerBound}}|-- Validation Window --|{{.UpperBound}}| ({{.UpperBound}}:HEAD)

`

var ParsedWindowTemplate = *template.Must(template.New("window").Parse(WindowTemplate))

type ValidationWindow struct {
	LowerBound int64
	UpperBound int64
}

func (window ValidationWindow) Size() int {
	return int(window.UpperBound - window.LowerBound)
}

func MakeValidationWindow(blockchain core.BlockChain, windowSize int) ValidationWindow {
	upperBound := blockchain.LastBlock().Int64()
	lowerBound := upperBound - int64(windowSize)
	return ValidationWindow{lowerBound, upperBound}
}

func MakeRange(min, max int64) []int64 {
	a := make([]int64, max-min+1)
	for i := range a {
		a[i] = min + int64(i)
	}
	return a
}

func (window ValidationWindow) Log(out io.Writer) {
	ParsedWindowTemplate.Execute(out, window)
}
