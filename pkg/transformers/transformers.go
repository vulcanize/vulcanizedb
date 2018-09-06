// Copyright 2018 Vulcanize
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package transformers

import (
	"github.com/vulcanize/vulcanizedb/pkg/transformers/bite"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/flip_kick"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/frob"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/pit_file"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/pit_file/debt_ceiling"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/pit_file/ilk"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/pit_file/stability_fee"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/price_feeds"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/tend"
)

func TransformerInitializers() []shared.TransformerInitializer {
	flipKickConfig := flip_kick.FlipKickConfig
	flipKickTransformerInitializer := flip_kick.FlipKickTransformerInitializer{Config: flipKickConfig}
	frobConfig := frob.FrobConfig
	frobTransformerInitializer := frob.FrobTransformerInitializer{Config: frobConfig}
	pitFileConfig := pit_file.PitFileConfig
	pitFileDebtCeilingTransformerInitializer := debt_ceiling.PitFileDebtCeilingTransformerInitializer{Config: pitFileConfig}
	pitFileIlkTransformerInitializer := ilk.PitFileIlkTransformerInitializer{Config: pitFileConfig}
	pitFileStabilityFeeTransformerInitializer := stability_fee.PitFileStabilityFeeTransformerInitializer{Config: pitFileConfig}
	priceFeedConfig := price_feeds.PriceFeedConfig
	priceFeedTransformerInitializer := price_feeds.PriceFeedTransformerInitializer{Config: priceFeedConfig}
	tendConfig := tend.TendConfig
	tendTransformerInitializer := tend.TendTransformerInitializer{Config: tendConfig}
	biteTransformerInitializer := bite.BiteTransformerInitializer{Config: bite.BiteConfig}

	return []shared.TransformerInitializer{
		biteTransformerInitializer.NewBiteTransformer,
		flipKickTransformerInitializer.NewFlipKickTransformer,
		frobTransformerInitializer.NewFrobTransformer,
		pitFileDebtCeilingTransformerInitializer.NewPitFileDebtCeilingTransformer,
		pitFileIlkTransformerInitializer.NewPitFileIlkTransformer,
		pitFileStabilityFeeTransformerInitializer.NewPitFileStabilityFeeTransformer,
		priceFeedTransformerInitializer.NewPriceFeedTransformer,
		tendTransformerInitializer.NewTendTransformer,
	}
}
