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
	"github.com/vulcanize/vulcanizedb/pkg/transformers/cat_file"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/cat_file/chop_lump"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/cat_file/flip"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/cat_file/pit_vow"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/deal"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/dent"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/drip_drip"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/drip_file"
	ilk2 "github.com/vulcanize/vulcanizedb/pkg/transformers/drip_file/ilk"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/drip_file/repo"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/drip_file/vow"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/flip_kick"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/flop_kick"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/frob"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/pit_file"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/pit_file/debt_ceiling"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/pit_file/ilk"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/pit_file/stability_fee"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/price_feeds"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/tend"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/vat_grab"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/vat_init"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/vat_toll"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/vat_tune"
)

var (
	BiteTransformerInitializer                = bite.BiteTransformerInitializer{Config: bite.BiteConfig}.NewBiteTransformer
	catFileConfig                             = cat_file.CatFileConfig
	CatFileChopLumpTransformerInitializer     = chop_lump.CatFileChopLumpTransformerInitializer{Config: catFileConfig}.NewCatFileChopLumpTransformer
	CatFileFlipTransformerInitializer         = flip.CatFileFlipTransformerInitializer{Config: catFileConfig}.NewCatFileFlipTransformer
	CatFilePitVowTransformerInitializer       = pit_vow.CatFilePitVowTransformerInitializer{Config: catFileConfig}.NewCatFilePitVowTransformer
	DealTransformerInitializer                = deal.DealTransformerInitializer{Config: deal.Config}.NewDealTransformer
	DentTransformerInitializer                = dent.DentTransformerInitializer{Config: dent.DentConfig}.NewDentTransformer
	DripDripTransformerInitializer            = drip_drip.DripDripTransformerInitializer{Config: drip_drip.DripDripConfig}.NewDripDripTransformer
	dripFileConfig                            = drip_file.DripFileConfig
	DripFileIlkTransformerInitializer         = ilk2.DripFileIlkTransformerInitializer{Config: dripFileConfig}.NewDripFileIlkTransformer
	DripFileRepoTransformerInitializer        = repo.DripFileRepoTransformerInitializer{Config: dripFileConfig}.NewDripFileRepoTransformer
	DripFileVowTransfromerInitializer         = vow.DripFileVowTransformerInitializer{Config: dripFileConfig}.NewDripFileVowTransformer
	FlipKickTransformerInitializer            = flip_kick.FlipKickTransformerInitializer{Config: flip_kick.FlipKickConfig}.NewFlipKickTransformer
	FlopKickTransformerInitializer            = flop_kick.FlopKickTransformerInitializer{Config: flop_kick.Config}.NewFlopKickTransformer
	FrobTransformerInitializer                = frob.FrobTransformerInitializer{Config: frob.FrobConfig}.NewFrobTransformer
	pitFileConfig                             = pit_file.PitFileConfig
	PitFileDebtCeilingTransformerInitializer  = debt_ceiling.PitFileDebtCeilingTransformerInitializer{Config: pitFileConfig}.NewPitFileDebtCeilingTransformer
	PitFileIlkTransformerInitializer          = ilk.PitFileIlkTransformerInitializer{Config: pitFileConfig}.NewPitFileIlkTransformer
	PitFileStabilityFeeTransformerInitializer = stability_fee.PitFileStabilityFeeTransformerInitializer{Config: pitFileConfig}.NewPitFileStabilityFeeTransformer
	PriceFeedTransformerInitializer           = price_feeds.PriceFeedTransformerInitializer{Config: price_feeds.PriceFeedConfig}.NewPriceFeedTransformer
	TendTransformerInitializer                = tend.TendTransformerInitializer{Config: tend.TendConfig}.NewTendTransformer
	VatGrabTransformerInitializer             = vat_grab.VatGrabTransformerInitializer{Config: vat_grab.VatGrabConfig}.NewVatGrabTransformer
	VatInitTransformerInitializer             = vat_init.VatInitTransformerInitializer{Config: vat_init.VatInitConfig}.NewVatInitTransformer
	VatTollTransformerInitializer             = vat_toll.VatTollTransformerInitializer{Config: vat_toll.VatTollConfig}.NewVatTollTransformer
	VatTuneTransformerInitializer             = vat_tune.VatTuneTransformerInitializer{Config: vat_tune.VatTuneConfig}.NewVatTuneTransformer
)

func TransformerInitializers() []shared.TransformerInitializer {
	return []shared.TransformerInitializer{
		BiteTransformerInitializer,
		CatFileChopLumpTransformerInitializer,
		CatFileFlipTransformerInitializer,
		CatFilePitVowTransformerInitializer,
		DealTransformerInitializer,
		DentTransformerInitializer,
		DripDripTransformerInitializer,
		DripFileIlkTransformerInitializer,
		DripFileVowTransfromerInitializer,
		DripFileRepoTransformerInitializer,
		FlipKickTransformerInitializer,
		FlopKickTransformerInitializer,
		FrobTransformerInitializer,
		PitFileDebtCeilingTransformerInitializer,
		PitFileIlkTransformerInitializer,
		PitFileStabilityFeeTransformerInitializer,
		PriceFeedTransformerInitializer,
		TendTransformerInitializer,
		VatGrabTransformerInitializer,
		VatInitTransformerInitializer,
		VatTollTransformerInitializer,
		VatTuneTransformerInitializer,
	}
}
