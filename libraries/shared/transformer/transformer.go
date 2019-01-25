// Auto-gen this code for different transformer interfaces/configs
// based on config file to allow for more modularity

package transformer

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/vulcanizedb/pkg/core"

	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared/constants"
)

type Transformer interface {
	Execute(logs []types.Log, header core.Header, recheckHeaders constants.TransformerExecution) error
	GetConfig() TransformerConfig
}

type TransformerInitializer func(db *postgres.DB) Transformer

type TransformerConfig struct {
	TransformerName     string
	ContractAddresses   []string
	ContractAbi         string
	Topic               string
	StartingBlockNumber int64
	EndingBlockNumber   int64 // Set -1 for indefinite transformer
}

func HexToInt64(byteString string) int64 {
	value := common.HexToHash(byteString)
	return value.Big().Int64()
}

func HexToString(byteString string) string {
	value := common.HexToHash(byteString)
	return value.Big().String()
}

func HexStringsToAddresses(strings []string) (addresses []common.Address) {
	for _, hexString := range strings {
		addresses = append(addresses, common.HexToAddress(hexString))
	}
	return
}
