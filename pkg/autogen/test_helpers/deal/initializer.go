package deal

import (
	"github.com/vulcanize/vulcanizedb/libraries/shared/transformer"
	"github.com/vulcanize/vulcanizedb/pkg/transformers"
)

var TransformerInitializer transformer.TransformerInitializer = transformers.GetDealTransformer().NewLogNoteTransformer
