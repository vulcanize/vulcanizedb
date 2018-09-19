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
	"github.com/vulcanize/vulcanizedb/pkg/transformers/deal"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/dent"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/drip_drip"
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

var (
	BiteTransformerInitializer                = bite.BiteTransformerInitializer{Config: bite.BiteConfig}.NewBiteTransformer
	DealTransformerInitializer                = deal.DealTransformerInitializer{Config: deal.Config}.NewDealTransformer
	DentTransformerInitializer                = dent.DentTransformerInitializer{Config: dent.DentConfig}.NewDentTransformer
	DripDripTransformerInitializer            = drip_drip.DripDripTransformerInitializer{Config: drip_drip.DripDripConfig}.NewDripDripTransformer
	dripFileConfig                            = drip_file.DripFileConfig
	DripFileIlkTransformerInitializer         = ilk2.DripFileIlkTransformerInitializer{Config: dripFileConfig}.NewDripFileIlkTransformer
	DripFileRepoTransformerInitializer        = repo.DripFileRepoTransformerInitializer{Config: dripFileConfig}.NewDripFileRepoTransformer
	FlipKickTransformerInitializer            = flip_kick.FlipKickTransformerInitializer{Config: flip_kick.FlipKickConfig}.NewFlipKickTransformer
	FrobTransformerInitializer                = frob.FrobTransformerInitializer{Config: frob.FrobConfig}.NewFrobTransformer
	pitFileConfig                             = pit_file.PitFileConfig
	PitFileDebtCeilingTransformerInitializer  = debt_ceiling.PitFileDebtCeilingTransformerInitializer{Config: pitFileConfig}.NewPitFileDebtCeilingTransformer
	PitFileIlkTransformerInitializer          = ilk.PitFileIlkTransformerInitializer{Config: pitFileConfig}.NewPitFileIlkTransformer
	PitFileStabilityFeeTransformerInitializer = stability_fee.PitFileStabilityFeeTransformerInitializer{Config: pitFileConfig}.NewPitFileStabilityFeeTransformer
	PriceFeedTransformerInitializer           = price_feeds.PriceFeedTransformerInitializer{Config: price_feeds.PriceFeedConfig}.NewPriceFeedTransformer
	TendTransformerInitializer                = tend.TendTransformerInitializer{Config: tend.TendConfig}.NewTendTransformer
	VatInitTransformerInitializer             = vat_init.VatInitTransformerInitializer{Config: vat_init.VatInitConfig}.NewVatInitTransformer
)

func TransformerInitializers() []shared.TransformerInitializer {
	return []shared.TransformerInitializer{
		BiteTransformerInitializer,
		DealTransformerInitializer,
		DentTransformerInitializer,
		DripDripTransformerInitializer,
		DripFileIlkTransformerInitializer,
		DripFileRepoTransformerInitializer,
		FlipKickTransformerInitializer,
		FrobTransformerInitializer,
		PitFileDebtCeilingTransformerInitializer,
		PitFileIlkTransformerInitializer,
		PitFileStabilityFeeTransformerInitializer,
		PriceFeedTransformerInitializer,
		TendTransformerInitializer,
		VatInitTransformerInitializer,
	}
}
