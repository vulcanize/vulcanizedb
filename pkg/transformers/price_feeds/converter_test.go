// Copyright © 2018 Vulcanize
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package price_feeds_test

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/transformers/price_feeds"
)

var _ = Describe("Price feed converter", func() {
	It("converts a log to a price feed model", func() {
		medianizerAddress := common.HexToAddress("0x99041f808d598b782d5a3e498681c2452a31da08")
		blockNumber := uint64(6147230)
		txIndex := uint(119)
		// https://etherscan.io/tx/0xa51a50a2adbfba4e2ab3d72dfd67a21c769f1bc8d2b180663a15500a56cde58f
		log := types.Log{
			Address:     medianizerAddress,
			Topics:      []common.Hash{common.HexToHash(price_feeds.LogValueTopic0)},
			Data:        common.FromHex("00000000000000000000000000000000000000000000001486f658319fb0c100"),
			BlockNumber: blockNumber,
			TxHash:      common.HexToHash("0xa51a50a2adbfba4e2ab3d72dfd67a21c769f1bc8d2b180663a15500a56cde58f"),
			TxIndex:     txIndex,
			BlockHash:   common.HexToHash("0x27ecebbf69eefa3bb3cf65f472322a80ff4946653a50a2171dc605f49829467d"),
			Index:       0,
			Removed:     false,
		}
		converter := price_feeds.PriceFeedConverter{}
		headerID := int64(123)

		model := converter.ToModel(log, headerID)

		expectedModel := price_feeds.PriceFeedModel{
			BlockNumber:       blockNumber,
			HeaderID:          headerID,
			MedianizerAddress: medianizerAddress[:],
			UsdValue:          "378.6599388897",
			TransactionIndex:  txIndex,
		}
		Expect(model).To(Equal(expectedModel))
	})
})
