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
	"github.com/vulcanize/vulcanizedb/pkg/transformers/flap_kick"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/flip_kick"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/flop_kick"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/frob"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/pit_file/debt_ceiling"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/pit_file/ilk"
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
	"github.com/vulcanize/vulcanizedb/pkg/transformers/vow_flog"
)

// Custom event transformers
var (
	BiteTransformer = factories.Transformer{
		Config:     bite.BiteConfig,
		Converter:  &bite.BiteConverter{},
		Repository: &bite.BiteRepository{},
	}

	FlapKickTransformer = factories.Transformer{
		Config:     flap_kick.FlapKickConfig,
		Converter:  &flap_kick.FlapKickConverter{},
		Repository: &flap_kick.FlapKickRepository{},
	}

	FlipKickTransformer = factories.Transformer{
		Config:     flip_kick.FlipKickConfig,
		Converter:  &flip_kick.FlipKickConverter{},
		Repository: &flip_kick.FlipKickRepository{},
	}

	FrobTransformer = factories.Transformer{
		Config:     frob.FrobConfig,
		Converter:  &frob.FrobConverter{},
		Repository: &frob.FrobRepository{},
	}

	FlopKickTransformer = factories.Transformer{
		Config:     flop_kick.Config,
		Converter:  &flop_kick.FlopKickConverter{},
		Repository: &flop_kick.FlopKickRepository{},
	}

	customEventTransformers = []factories.Transformer{
		BiteTransformer,
		FlapKickTransformer,
		FlipKickTransformer,
		FrobTransformer,
		FlopKickTransformer,
	}
)

// LogNote transformers
var (
	CatFileChopLumpTransformer = factories.LogNoteTransformer{
		Config:     chop_lump.CatFileChopLumpConfig,
		Converter:  &chop_lump.CatFileChopLumpConverter{},
		Repository: &chop_lump.CatFileChopLumpRepository{},
	}

	CatFileFlipTransformer = factories.LogNoteTransformer{
		Config:     flip.CatFileFlipConfig,
		Converter:  &flip.CatFileFlipConverter{},
		Repository: &flip.CatFileFlipRepository{},
	}

	CatFilePitVowTransformer = factories.LogNoteTransformer{
		Config:     pit_vow.CatFilePitVowConfig,
		Converter:  &pit_vow.CatFilePitVowConverter{},
		Repository: &pit_vow.CatFilePitVowRepository{},
	}

	DealTransformer = factories.LogNoteTransformer{
		Config:     deal.DealConfig,
		Converter:  &deal.DealConverter{},
		Repository: &deal.DealRepository{},
	}

	DentTransformer = factories.LogNoteTransformer{
		Config:     dent.DentConfig,
		Converter:  &dent.DentConverter{},
		Repository: &dent.DentRepository{},
	}

	DripDripTransformer = factories.LogNoteTransformer{
		Config:     drip_drip.DripDripConfig,
		Converter:  &drip_drip.DripDripConverter{},
		Repository: &drip_drip.DripDripRepository{},
	}

	DripFileIlkTransformer = factories.LogNoteTransformer{
		Config:     ilk2.DripFileIlkConfig,
		Converter:  &ilk2.DripFileIlkConverter{},
		Repository: &ilk2.DripFileIlkRepository{},
	}

	DripFileRepoTransformer = factories.LogNoteTransformer{
		Config:     repo.DripFileRepoConfig,
		Converter:  &repo.DripFileRepoConverter{},
		Repository: &repo.DripFileRepoRepository{},
	}

	DripFileVowTransfromer = factories.LogNoteTransformer{
		Config:     vow.DripFileVowConfig,
		Converter:  &vow.DripFileVowConverter{},
		Repository: &vow.DripFileVowRepository{},
	}

	FlogTransformer = factories.LogNoteTransformer{
		Config:     vow_flog.VowFlogConfig,
		Converter:  &vow_flog.VowFlogConverter{},
		Repository: &vow_flog.VowFlogRepository{},
	}

	PitFileDebtCeilingTransformer = factories.LogNoteTransformer{
		Config:     debt_ceiling.DebtCeilingFileConfig,
		Converter:  &debt_ceiling.PitFileDebtCeilingConverter{},
		Repository: &debt_ceiling.PitFileDebtCeilingRepository{},
	}

	PitFileIlkTransformer = factories.LogNoteTransformer{
		Config:     ilk.IlkFileConfig,
		Converter:  &ilk.PitFileIlkConverter{},
		Repository: &ilk.PitFileIlkRepository{},
	}

	PriceFeedTransformer = factories.LogNoteTransformer{
		Config:     price_feeds.PriceFeedConfig,
		Converter:  &price_feeds.PriceFeedConverter{},
		Repository: &price_feeds.PriceFeedRepository{},
	}

	TendTransformer = factories.LogNoteTransformer{
		Config:     tend.TendConfig,
		Converter:  &tend.TendConverter{},
		Repository: &tend.TendRepository{},
	}

	VatInitTransformer = factories.LogNoteTransformer{
		Config:     vat_init.VatInitConfig,
		Converter:  &vat_init.VatInitConverter{},
		Repository: &vat_init.VatInitRepository{},
	}

	VatGrabTransformer = factories.LogNoteTransformer{
		Config:     vat_grab.VatGrabConfig,
		Converter:  &vat_grab.VatGrabConverter{},
		Repository: &vat_grab.VatGrabRepository{},
	}

	VatFoldTransformer = factories.LogNoteTransformer{
		Config:     vat_fold.VatFoldConfig,
		Converter:  &vat_fold.VatFoldConverter{},
		Repository: &vat_fold.VatFoldRepository{},
	}

	VatHealTransformer = factories.LogNoteTransformer{
		Config:     vat_heal.VatHealConfig,
		Converter:  &vat_heal.VatHealConverter{},
		Repository: &vat_heal.VatHealRepository{},
	}

	VatMoveTransformer = factories.LogNoteTransformer{
		Config:     vat_move.VatMoveConfig,
		Converter:  &vat_move.VatMoveConverter{},
		Repository: &vat_move.VatMoveRepository{},
	}

	VatSlipTransformer = factories.LogNoteTransformer{
		Config:     vat_slip.VatSlipConfig,
		Converter:  &vat_slip.VatSlipConverter{},
		Repository: &vat_slip.VatSlipRepository{},
	}

	VatTollTransformer = factories.LogNoteTransformer{
		Config:     vat_toll.VatTollConfig,
		Converter:  &vat_toll.VatTollConverter{},
		Repository: &vat_toll.VatTollRepository{},
	}

	VatTuneTransformer = factories.LogNoteTransformer{
		Config:     vat_tune.VatTuneConfig,
		Converter:  &vat_tune.VatTuneConverter{},
		Repository: &vat_tune.VatTuneRepository{},
	}

	VatFluxTransformer = factories.LogNoteTransformer{
		Config:     vat_flux.VatFluxConfig,
		Converter:  &vat_flux.VatFluxConverter{},
		Repository: &vat_flux.VatFluxRepository{},
	}

	logNoteTransformers = []factories.LogNoteTransformer{
		CatFileChopLumpTransformer,
		CatFileFlipTransformer,
		CatFilePitVowTransformer,
		DealTransformer,
		DentTransformer,
		DripDripTransformer,
		DripFileIlkTransformer,
		DripFileRepoTransformer,
		DripFileVowTransfromer,
		FlogTransformer,
		PitFileDebtCeilingTransformer,
		PitFileIlkTransformer,
		PriceFeedTransformer,
		TendTransformer,
		VatInitTransformer,
		VatGrabTransformer,
		VatFoldTransformer,
		VatHealTransformer,
		VatMoveTransformer,
		VatSlipTransformer,
		VatTollTransformer,
		VatTuneTransformer,
		VatFluxTransformer,
	}
)

// `TransformerInitializers` returns a list of functions, that given a db pointer
// will return a `shared.Transformer`
func TransformerInitializers() (initializers []shared.TransformerInitializer) {
	for _, transformer := range logNoteTransformers {
		initializers = append(initializers, transformer.NewLogNoteTransformer)
	}

	for _, transformer := range customEventTransformers {
		initializers = append(initializers, transformer.NewTransformer)
	}
	return
}
