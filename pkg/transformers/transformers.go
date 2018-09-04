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
	"github.com/vulcanize/vulcanizedb/pkg/transformers/dent"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/drip_file"
	ilk2 "github.com/vulcanize/vulcanizedb/pkg/transformers/drip_file/ilk"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/drip_file/repo"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/flip_kick"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/frob"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/pit_file"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/pit_file/debt_ceiling"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/pit_file/ilk"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/pit_file/stability_fee"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/price_feeds"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/tend"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/vat_init"
)

func TransformerInitializers() []shared.TransformerInitializer {
	biteTransformerInitializer := bite.BiteTransformerInitializer{Config: bite.BiteConfig}
	dentTransformerInitializer := dent.DentTransformerInitializer{Config: dent.DentConfig}
	flipKickTransformerInitializer := flip_kick.FlipKickTransformerInitializer{Config: flip_kick.FlipKickConfig}
	frobTransformerInitializer := frob.FrobTransformerInitializer{Config: frob.FrobConfig}
	dripFileConfig := drip_file.DripFileConfig
	dripFileIlkTransformerInitializer := ilk2.DripFileIlkTransformerInitializer{Config: dripFileConfig}
	dripFileRepoTransformerInitializer := repo.DripFileRepoTransformerInitializer{Config: dripFileConfig}
	pitFileConfig := pit_file.PitFileConfig
	pitFileDebtCeilingTransformerInitializer := debt_ceiling.PitFileDebtCeilingTransformerInitializer{Config: pitFileConfig}
	pitFileIlkTransformerInitializer := ilk.PitFileIlkTransformerInitializer{Config: pitFileConfig}
	pitFileStabilityFeeTransformerInitializer := stability_fee.PitFileStabilityFeeTransformerInitializer{Config: pitFileConfig}
	priceFeedTransformerInitializer := price_feeds.PriceFeedTransformerInitializer{Config: price_feeds.PriceFeedConfig}
	tendTransformerInitializer := tend.TendTransformerInitializer{Config: tend.TendConfig}
	vatInitConfig := vat_init.VatInitConfig
	vatInitTransformerInitializer := vat_init.VatInitTransformerInitializer{Config: vatInitConfig}

	return []shared.TransformerInitializer{
		biteTransformerInitializer.NewBiteTransformer,
		dentTransformerInitializer.NewDentTransformer,
		dripFileIlkTransformerInitializer.NewDripFileIlkTransformer,
		dripFileRepoTransformerInitializer.NewDripFileRepoTransformer,
		flipKickTransformerInitializer.NewFlipKickTransformer,
		frobTransformerInitializer.NewFrobTransformer,
		pitFileDebtCeilingTransformerInitializer.NewPitFileDebtCeilingTransformer,
		pitFileIlkTransformerInitializer.NewPitFileIlkTransformer,
		pitFileStabilityFeeTransformerInitializer.NewPitFileStabilityFeeTransformer,
		priceFeedTransformerInitializer.NewPriceFeedTransformer,
		tendTransformerInitializer.NewTendTransformer,
		vatInitTransformerInitializer.NewVatInitTransformer,
	}
}
