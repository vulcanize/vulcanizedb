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
func GetBiteTransformer() factories.Transformer {
	return factories.Transformer{
		Config:     bite.GetBiteConfig(),
		Converter:  &bite.BiteConverter{},
		Repository: &bite.BiteRepository{},
	}
}

func GetFlapKickTransformer() factories.Transformer {
	return factories.Transformer{
		Config:     flap_kick.GetFlapKickConfig(),
		Converter:  &flap_kick.FlapKickConverter{},
		Repository: &flap_kick.FlapKickRepository{},
	}
}

func GetFlipKickTransformer() factories.Transformer {
	return factories.Transformer{
		Config:     flip_kick.GetFlipKickConfig(),
		Converter:  &flip_kick.FlipKickConverter{},
		Repository: &flip_kick.FlipKickRepository{},
	}
}

func GetFrobTransformer() factories.Transformer {
	return factories.Transformer{
		Config:     frob.GetFrobConfig(),
		Converter:  &frob.FrobConverter{},
		Repository: &frob.FrobRepository{},
	}
}

func GetFlopKickTransformer() factories.Transformer {
	return factories.Transformer{
		Config:     flop_kick.GetFlopKickConfig(),
		Converter:  &flop_kick.FlopKickConverter{},
		Repository: &flop_kick.FlopKickRepository{},
	}
}

func getCustomEventTransformers() []factories.Transformer {
	return []factories.Transformer{
		GetBiteTransformer(),
		GetFlapKickTransformer(),
		GetFlipKickTransformer(),
		GetFrobTransformer(),
		GetFlopKickTransformer(),
	}
}

// LogNote transformers
func GetCatFileChopLumpTransformer() factories.LogNoteTransformer {
	return factories.LogNoteTransformer{
		Config:     chop_lump.GetCatFileChopLumpConfig(),
		Converter:  &chop_lump.CatFileChopLumpConverter{},
		Repository: &chop_lump.CatFileChopLumpRepository{},
	}
}

func GetCatFileFlipTransformer() factories.LogNoteTransformer {
	return factories.LogNoteTransformer{
		Config:     flip.GetCatFileFlipConfig(),
		Converter:  &flip.CatFileFlipConverter{},
		Repository: &flip.CatFileFlipRepository{},
	}
}

func GetCatFilePitVowTransformer() factories.LogNoteTransformer {
	return factories.LogNoteTransformer{
		Config:     pit_vow.GetCatFilePitVowConfig(),
		Converter:  &pit_vow.CatFilePitVowConverter{},
		Repository: &pit_vow.CatFilePitVowRepository{},
	}
}

func GetDealTransformer() factories.LogNoteTransformer {
	return factories.LogNoteTransformer{
		Config:     deal.GetDealConfig(),
		Converter:  &deal.DealConverter{},
		Repository: &deal.DealRepository{},
	}
}

func GetDentTransformer() factories.LogNoteTransformer {
	return factories.LogNoteTransformer{
		Config:     dent.GetDentConfig(),
		Converter:  &dent.DentConverter{},
		Repository: &dent.DentRepository{},
	}
}

func GetDripDripTransformer() factories.LogNoteTransformer {
	return factories.LogNoteTransformer{
		Config:     drip_drip.GetDripDripConfig(),
		Converter:  &drip_drip.DripDripConverter{},
		Repository: &drip_drip.DripDripRepository{},
	}
}

func GetDripFileIlkTransformer() factories.LogNoteTransformer {
	return factories.LogNoteTransformer{
		Config:     ilk2.GetDripFileIlkConfig(),
		Converter:  &ilk2.DripFileIlkConverter{},
		Repository: &ilk2.DripFileIlkRepository{},
	}
}

func GetDripFileRepoTransformer() factories.LogNoteTransformer {
	return factories.LogNoteTransformer{
		Config:     repo.GetDripFileRepoConfig(),
		Converter:  &repo.DripFileRepoConverter{},
		Repository: &repo.DripFileRepoRepository{},
	}
}

func GetDripFileVowTransformer() factories.LogNoteTransformer {
	return factories.LogNoteTransformer{
		Config:     vow.GetDripFileVowConfig(),
		Converter:  &vow.DripFileVowConverter{},
		Repository: &vow.DripFileVowRepository{},
	}
}

func GetFlogTransformer() factories.LogNoteTransformer {
	return factories.LogNoteTransformer{
		Config:     vow_flog.GetVowFlogConfig(),
		Converter:  &vow_flog.VowFlogConverter{},
		Repository: &vow_flog.VowFlogRepository{},
	}
}

func GetPitFileDebtCeilingTransformer() factories.LogNoteTransformer {
	return factories.LogNoteTransformer{
		Config:     debt_ceiling.GetDebtCeilingFileConfig(),
		Converter:  &debt_ceiling.PitFileDebtCeilingConverter{},
		Repository: &debt_ceiling.PitFileDebtCeilingRepository{},
	}
}

func GetPitFileIlkTransformer() factories.LogNoteTransformer {
	return factories.LogNoteTransformer{
		Config:     ilk.GetIlkFileConfig(),
		Converter:  &ilk.PitFileIlkConverter{},
		Repository: &ilk.PitFileIlkRepository{},
	}
}

func GetPriceFeedTransformer() factories.LogNoteTransformer {
	return factories.LogNoteTransformer{
		Config:     price_feeds.GetPriceFeedConfig(),
		Converter:  &price_feeds.PriceFeedConverter{},
		Repository: &price_feeds.PriceFeedRepository{},
	}
}

func GetTendTransformer() factories.LogNoteTransformer {
	return factories.LogNoteTransformer{
		Config:     tend.GetTendConfig(),
		Converter:  &tend.TendConverter{},
		Repository: &tend.TendRepository{},
	}
}

func GetVatInitTransformer() factories.LogNoteTransformer {
	return factories.LogNoteTransformer{
		Config:     vat_init.GetVatInitConfig(),
		Converter:  &vat_init.VatInitConverter{},
		Repository: &vat_init.VatInitRepository{},
	}
}

func GetVatGrabTransformer() factories.LogNoteTransformer {
	return factories.LogNoteTransformer{
		Config:     vat_grab.GetVatGrabConfig(),
		Converter:  &vat_grab.VatGrabConverter{},
		Repository: &vat_grab.VatGrabRepository{},
	}
}

func GetVatFoldTransformer() factories.LogNoteTransformer {
	return factories.LogNoteTransformer{
		Config:     vat_fold.GetVatFoldConfig(),
		Converter:  &vat_fold.VatFoldConverter{},
		Repository: &vat_fold.VatFoldRepository{},
	}
}

func GetVatHealTransformer() factories.LogNoteTransformer {
	return factories.LogNoteTransformer{
		Config:     vat_heal.GetVatHealConfig(),
		Converter:  &vat_heal.VatHealConverter{},
		Repository: &vat_heal.VatHealRepository{},
	}
}

func GetVatMoveTransformer() factories.LogNoteTransformer {
	return factories.LogNoteTransformer{
		Config:     vat_move.GetVatMoveConfig(),
		Converter:  &vat_move.VatMoveConverter{},
		Repository: &vat_move.VatMoveRepository{},
	}
}

func GetVatSlipTransformer() factories.LogNoteTransformer {
	return factories.LogNoteTransformer{
		Config:     vat_slip.GetVatSlipConfig(),
		Converter:  &vat_slip.VatSlipConverter{},
		Repository: &vat_slip.VatSlipRepository{},
	}
}

func GetVatTollTransformer() factories.LogNoteTransformer {
	return factories.LogNoteTransformer{
		Config:     vat_toll.GetVatTollConfig(),
		Converter:  &vat_toll.VatTollConverter{},
		Repository: &vat_toll.VatTollRepository{},
	}
}

func GetVatTuneTransformer() factories.LogNoteTransformer {
	return factories.LogNoteTransformer{
		Config:     vat_tune.GetVatTuneConfig(),
		Converter:  &vat_tune.VatTuneConverter{},
		Repository: &vat_tune.VatTuneRepository{},
	}
}

func GetVatFluxTransformer() factories.LogNoteTransformer {
	return factories.LogNoteTransformer{
		Config:     vat_flux.GetVatFluxConfig(),
		Converter:  &vat_flux.VatFluxConverter{},
		Repository: &vat_flux.VatFluxRepository{},
	}
}

func getLogNoteTransformers() []factories.LogNoteTransformer {
	return []factories.LogNoteTransformer{
		GetCatFileChopLumpTransformer(),
		GetCatFileFlipTransformer(),
		GetCatFilePitVowTransformer(),
		GetDealTransformer(),
		GetDentTransformer(),
		GetDripDripTransformer(),
		GetDripFileIlkTransformer(),
		GetDripFileRepoTransformer(),
		GetDripFileVowTransformer(),
		GetFlogTransformer(),
		GetPitFileDebtCeilingTransformer(),
		GetPitFileIlkTransformer(),
		GetPriceFeedTransformer(),
		GetTendTransformer(),
		GetVatInitTransformer(),
		GetVatGrabTransformer(),
		GetVatFoldTransformer(),
		GetVatHealTransformer(),
		GetVatMoveTransformer(),
		GetVatSlipTransformer(),
		GetVatTollTransformer(),
		GetVatTuneTransformer(),
		GetVatFluxTransformer(),
	}
}

// `TransformerInitializers` returns a list of functions, that given a db pointer
// will return a `shared.Transformer`
func TransformerInitializers() (initializers []shared.TransformerInitializer) {
	for _, transformer := range getLogNoteTransformers() {
		initializers = append(initializers, transformer.NewLogNoteTransformer)
	}

	for _, transformer := range getCustomEventTransformers() {
		initializers = append(initializers, transformer.NewTransformer)
	}
	return
}
