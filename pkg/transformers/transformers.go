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
	"github.com/vulcanize/vulcanizedb/pkg/transformers/cat_file/chop_lump"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/cat_file/flip"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/cat_file/pit_vow"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/deal"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/dent"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/drip_drip"
	ilk2 "github.com/vulcanize/vulcanizedb/pkg/transformers/drip_file/ilk"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/drip_file/repo"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/drip_file/vow"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/factories"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/flip_kick"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/flog"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/flop_kick"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/frob"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/pit_file/debt_ceiling"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/pit_file/ilk"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/pit_file/stability_fee"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/price_feeds"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/tend"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/vat_flux"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/vat_fold"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/vat_grab"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/vat_heal"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/vat_init"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/vat_move"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/vat_slip"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/vat_toll"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/vat_tune"
)

var (
	BiteTransformerInitializer = factories.Transformer{
		Config:     bite.BiteConfig,
		Converter:  &bite.BiteConverter{},
		Repository: &bite.BiteRepository{},
		Fetcher:    &shared.Fetcher{},
	}.NewTransformer

	CatFileChopLumpTransformerInitializer = factories.LogNoteTransformer{
		Config:     chop_lump.CatFileChopLumpConfig,
		Converter:  &chop_lump.CatFileChopLumpConverter{},
		Repository: &chop_lump.CatFileChopLumpRepository{},
		Fetcher:    &shared.Fetcher{},
	}.NewLogNoteTransformer

	CatFileFlipTransformerInitializer = factories.LogNoteTransformer{
		Config:     flip.CatFileFlipConfig,
		Converter:  &flip.CatFileFlipConverter{},
		Repository: &flip.CatFileFlipRepository{},
		Fetcher:    &shared.Fetcher{},
	}.NewLogNoteTransformer

	CatFilePitVowTransformerInitializer = factories.LogNoteTransformer{
		Config:     pit_vow.CatFilePitVowConfig,
		Converter:  &pit_vow.CatFilePitVowConverter{},
		Repository: &pit_vow.CatFilePitVowRepository{},
		Fetcher:    &shared.Fetcher{},
	}.NewLogNoteTransformer

	DealTransformerInitializer = factories.LogNoteTransformer{
		Config:     deal.DealConfig,
		Converter:  &deal.DealConverter{},
		Repository: &deal.DealRepository{},
		Fetcher:    &shared.Fetcher{},
	}.NewLogNoteTransformer

	DentTransformerInitializer = factories.Transformer{
		Config:     dent.DentConfig,
		Converter:  &dent.DentConverter{},
		Repository: &dent.DentRepository{},
		Fetcher:    &shared.Fetcher{},
	}.NewTransformer

	DripDripTransformerInitializer    = drip_drip.DripDripTransformerInitializer{Config: drip_drip.DripDripConfig}.NewDripDripTransformer

	DripFileIlkTransformerInitializer = factories.LogNoteTransformer{
		Config:     ilk2.DripFileIlkConfig,
		Converter:  &ilk2.DripFileIlkConverter{},
		Repository: &ilk2.DripFileIlkRepository{},
		Fetcher:    &shared.Fetcher{},
	}.NewLogNoteTransformer

	DripFileRepoTransformerInitializer = factories.LogNoteTransformer{
		Config:     repo.DripFileRepoConfig,
		Converter:  &repo.DripFileRepoConverter{},
		Repository: &repo.DripFileRepoRepository{},
		Fetcher:    &shared.Fetcher{},
	}.NewLogNoteTransformer

	DripFileVowTransfromerInitializer = factories.LogNoteTransformer{
		Config:     vow.DripFileVowConfig,
		Converter:  &vow.DripFileVowConverter{},
		Repository: &vow.DripFileVowRepository{},
		Fetcher:    &shared.Fetcher{},
	}.NewLogNoteTransformer

	FlipKickTransformerInitializer = flip_kick.FlipKickTransformerInitializer{Config: flip_kick.FlipKickConfig}.NewFlipKickTransformer
	FlogTransformerInitializer     = factories.LogNoteTransformer{
		Config:     flog.FlogConfig,
		Converter:  &flog.FlogConverter{},
		Repository: &flog.FlogRepository{},
		Fetcher:    &shared.Fetcher{},
	}.NewLogNoteTransformer

	FlopKickTransformerInitializer           = flop_kick.FlopKickTransformerInitializer{Config: flop_kick.Config}.NewFlopKickTransformer
	FrobTransformerInitializer               = frob.FrobTransformerInitializer{Config: frob.FrobConfig}.NewFrobTransformer
	PitFileDebtCeilingTransformerInitializer = factories.LogNoteTransformer{
		Config:     debt_ceiling.DebtCeilingFileConfig,
		Converter:  &debt_ceiling.PitFileDebtCeilingConverter{},
		Repository: &debt_ceiling.PitFileDebtCeilingRepository{},
		Fetcher:    &shared.Fetcher{},
	}.NewLogNoteTransformer

	PitFileIlkTransformerInitializer = factories.LogNoteTransformer{
		Config:     ilk.IlkFileConfig,
		Converter:  &ilk.PitFileIlkConverter{},
		Repository: &ilk.PitFileIlkRepository{},
		Fetcher:    &shared.Fetcher{},
	}.NewLogNoteTransformer

	PitFileStabilityFeeTransformerInitializer = factories.LogNoteTransformer{
		Config:     stability_fee.StabilityFeeFileConfig,
		Converter:  &stability_fee.PitFileStabilityFeeConverter{},
		Repository: &stability_fee.PitFileStabilityFeeRepository{},
		Fetcher:    &shared.Fetcher{},
	}.NewLogNoteTransformer

	PriceFeedTransformerInitializer = factories.LogNoteTransformer{
		Config:     price_feeds.PriceFeedConfig,
		Converter:  &price_feeds.PriceFeedConverter{},
		Repository: &price_feeds.PriceFeedRepository{},
		Fetcher:    &shared.Fetcher{},
	}.NewLogNoteTransformer

	TendTransformerInitializer = factories.LogNoteTransformer{
		Config:     tend.TendConfig,
		Converter:  &tend.TendConverter{},
		Repository: &tend.TendRepository{},
		Fetcher:    &shared.Fetcher{},
	}.NewLogNoteTransformer

	VatInitTransformerInitializer = factories.LogNoteTransformer{
		Config:     vat_init.VatInitConfig,
		Converter:  &vat_init.VatInitConverter{},
		Repository: &vat_init.VatInitRepository{},
		Fetcher:    &shared.Fetcher{},
	}.NewLogNoteTransformer

	VatGrabTransformerInitializer = factories.LogNoteTransformer{
		Config:     vat_grab.VatGrabConfig,
		Converter:  &vat_grab.VatGrabConverter{},
		Repository: &vat_grab.VatGrabRepository{},
		Fetcher:    &shared.Fetcher{},
	}.NewLogNoteTransformer

	VatFoldTransformerInitializer = factories.LogNoteTransformer{
		Config:     vat_fold.VatFoldConfig,
		Converter:  &vat_fold.VatFoldConverter{},
		Repository: &vat_fold.VatFoldRepository{},
		Fetcher:    &shared.Fetcher{},
	}.NewLogNoteTransformer

	VatHealTransformerInitializer = factories.LogNoteTransformer{
		Config:     vat_heal.VatHealConfig,
		Converter:  &vat_heal.VatHealConverter{},
		Repository: &vat_heal.VatHealRepository{},
		Fetcher:    &shared.Fetcher{},
	}.NewLogNoteTransformer

	VatMoveTransformerInitializer = factories.LogNoteTransformer{
		Config:     vat_move.VatMoveConfig,
		Converter:  &vat_move.VatMoveConverter{},
		Repository: &vat_move.VatMoveRepository{},
		Fetcher:    &shared.Fetcher{},
	}.NewLogNoteTransformer

	VatSlipTransformerInitializer = factories.LogNoteTransformer{
		Config:     vat_slip.VatSlipConfig,
		Converter:  &vat_slip.VatSlipConverter{},
		Repository: &vat_slip.VatSlipRepository{},
		Fetcher:    &shared.Fetcher{},
	}.NewLogNoteTransformer

	VatTollTransformerInitializer = factories.LogNoteTransformer{
		Config:     vat_toll.VatTollConfig,
		Converter:  &vat_toll.VatTollConverter{},
		Repository: &vat_toll.VatTollRepository{},
		Fetcher:    &shared.Fetcher{},
	}.NewLogNoteTransformer

	VatTuneTransformerInitializer = factories.LogNoteTransformer{
		Config:     vat_tune.VatTuneConfig,
		Converter:  &vat_tune.VatTuneConverter{},
		Repository: &vat_tune.VatTuneRepository{},
		Fetcher:    &shared.Fetcher{},
	}.NewLogNoteTransformer

	VatFluxTransformerInitializer = factories.LogNoteTransformer{
		Config:     vat_flux.VatFluxConfig,
		Converter:  &vat_flux.VatFluxConverter{},
		Repository: &vat_flux.VatFluxRepository{},
		Fetcher:    &shared.Fetcher{},
	}.NewLogNoteTransformer
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
		FlogTransformerInitializer,
		FlopKickTransformerInitializer,
		FrobTransformerInitializer,
		PitFileDebtCeilingTransformerInitializer,
		PitFileIlkTransformerInitializer,
		PitFileStabilityFeeTransformerInitializer,
		PriceFeedTransformerInitializer,
		TendTransformerInitializer,
		VatGrabTransformerInitializer,
		VatInitTransformerInitializer,
		VatMoveTransformerInitializer,
		VatHealTransformerInitializer,
		VatFoldTransformerInitializer,
		VatSlipTransformerInitializer,
		VatTollTransformerInitializer,
		VatTuneTransformerInitializer,
		VatFluxTransformerInitializer,
	}
}
