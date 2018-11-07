// VulcanizeDB
// Copyright © 2018 Vulcanize

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.

// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package common_test

import (
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	vulcCommon "github.com/vulcanize/vulcanizedb/pkg/geth/converters/common"
)

var _ = Describe("Conversion of GethLog to core.Log", func() {

	It("converts geth log to internal log format", func() {
		gethLog := types.Log{
			Address:     common.HexToAddress("0x448a5065aeBB8E423F0896E6c5D525C040f59af3"),
			BlockHash:   common.HexToHash("0x656c34545f90a730a19008c0e7a7cd4fb3895064b48d6d69761bd5abad681056"),
			BlockNumber: 2019236,
			Data:        hexutil.MustDecode("0x000000000000000000000000000000000000000000000001a055690d9db80000"),
			Index:       2,
			TxIndex:     3,
			TxHash:      common.HexToHash("0x3b198bfd5d2907285af009e9ae84a0ecd63677110d89d7e030251acb87f6487e"),
			Topics: []common.Hash{
				common.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"),
				common.HexToHash("0x00000000000000000000000080b2c9d7cbbf30a1b0fc8983c647d754c6525615"),
			},
		}

		expected := core.Log{
			Address:     strings.ToLower(gethLog.Address.Hex()),
			BlockNumber: int64(gethLog.BlockNumber),
			Data:        hexutil.Encode(gethLog.Data),
			TxHash:      gethLog.TxHash.Hex(),
			Index:       2,
			Topics: core.Topics{
				gethLog.Topics[0].Hex(),
				gethLog.Topics[1].Hex(),
			},
		}

		coreLog := vulcCommon.ToCoreLog(gethLog)

		Expect(coreLog.Address).To(Equal(expected.Address))
		Expect(coreLog.BlockNumber).To(Equal(expected.BlockNumber))
		Expect(coreLog.Data).To(Equal(expected.Data))
		Expect(coreLog.Index).To(Equal(expected.Index))
		Expect(coreLog.Topics[0]).To(Equal(expected.Topics[0]))
		Expect(coreLog.Topics[1]).To(Equal(expected.Topics[1]))
		Expect(coreLog.TxHash).To(Equal(expected.TxHash))
	})

	It("converts geth log array to array of internal logs", func() {
		gethLogOne := types.Log{
			Address:     common.HexToAddress("0xecf8f87f810ecf450940c9f60066b4a7a501d6a7"),
			BlockHash:   common.HexToHash("0x656c34545f90a730a19008c0e7a7cd4fb3895064b48d6d69761bd5abad681056"),
			BlockNumber: 2019236,
			Data:        hexutil.MustDecode("0x000000000000000000000000000000000000000000000001a055690d9db80000"),
			Index:       2,
			TxIndex:     3,
			TxHash:      common.HexToHash("0x3b198bfd5d2907285af009e9ae84a0ecd63677110d89d7e030251acb87f6487e"),
			Topics: []common.Hash{
				common.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"),
				common.HexToHash("0x00000000000000000000000080b2c9d7cbbf30a1b0fc8983c647d754c6525615"),
			},
		}

		gethLogTwo := types.Log{
			Address:     common.HexToAddress("0x123"),
			BlockHash:   common.HexToHash("0x576"),
			BlockNumber: 2019236,
			Data:        hexutil.MustDecode("0x0000000000000000000000000000000000000000000000000000000000000001"),
			Index:       3,
			TxIndex:     4,
			TxHash:      common.HexToHash("0x134"),
			Topics: []common.Hash{
				common.HexToHash("0xaaa"),
				common.HexToHash("0xbbb"),
			},
		}

		expectedOne := vulcCommon.ToCoreLog(gethLogOne)
		expectedTwo := vulcCommon.ToCoreLog(gethLogTwo)

		coreLogs := vulcCommon.ToCoreLogs([]types.Log{gethLogOne, gethLogTwo})

		Expect(len(coreLogs)).To(Equal(2))
		Expect(coreLogs[0]).To(Equal(expectedOne))
		Expect(coreLogs[1]).To(Equal(expectedTwo))
	})
})
