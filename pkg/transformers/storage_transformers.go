package transformers

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/factories/storage"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared/constants"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/storage_diffs/maker"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/storage_diffs/maker/cat"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/storage_diffs/maker/pit"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/storage_diffs/maker/vat"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/storage_diffs/maker/vow"
)

func GetCatStorageTransformer() storage.Transformer {
	return storage.Transformer{
		Address:    common.HexToAddress(constants.CatContractAddress()),
		Mappings:   &cat.CatMappings{StorageRepository: &maker.MakerStorageRepository{}},
		Repository: &cat.CatStorageRepository{},
	}
}

func GetPitStorageTransformer() storage.Transformer {
	return storage.Transformer{
		Address:    common.HexToAddress(constants.PitContractAddress()),
		Mappings:   &pit.PitMappings{StorageRepository: &maker.MakerStorageRepository{}},
		Repository: &pit.PitStorageRepository{},
	}
}

func GetVatStorageTransformer() storage.Transformer {
	return storage.Transformer{
		Address:    common.HexToAddress(constants.VatContractAddress()),
		Mappings:   &vat.VatMappings{StorageRepository: &maker.MakerStorageRepository{}},
		Repository: &vat.VatStorageRepository{},
	}
}

func GetVowStorageTransformer() storage.Transformer {
	return storage.Transformer{
		Address:    common.HexToAddress(constants.VowContractAddress()),
		Mappings:   &vow.VowMappings{StorageRepository: &maker.MakerStorageRepository{}},
		Repository: &vow.VowStorageRepository{},
	}
}
