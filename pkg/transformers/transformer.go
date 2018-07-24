package transformers

import (
	"github.com/vulcanize/vulcanizedb/pkg/core"
)

type Transformer interface {
	Execute(header core.Header, headerID int64) error
}
