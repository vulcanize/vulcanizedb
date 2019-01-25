package bite

import (
	"github.com/vulcanize/vulcanizedb/libraries/shared/transformer"
	"github.com/vulcanize/vulcanizedb/pkg/transformers"
)

var TransformerInitializer transformer.TransformerInitializer = transformers.GetBiteTransformer().NewTransformer
