// VulcanizeDB
// Copyright © 2019 Vulcanize

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

package ipfs_test

import (
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rlp"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/ipfs"
	"github.com/vulcanize/vulcanizedb/pkg/ipfs/mocks"
)

var _ = Describe("Converter", func() {
	Describe("Convert", func() {
		It("Converts mock statediff.Payloads into the expected IPLDPayloads", func() {
			converter := ipfs.NewPayloadConverter(params.MainnetChainConfig)
			converterPayload, err := converter.Convert(mocks.MockStateDiffPayload)
			Expect(err).ToNot(HaveOccurred())
			Expect(converterPayload.BlockNumber).To(Equal(mocks.BlockNumber))
			Expect(converterPayload.BlockHash).To(Equal(mocks.MockBlock.Hash()))
			Expect(converterPayload.StateNodes).To(Equal(mocks.MockStateNodes))
			Expect(converterPayload.StorageNodes).To(Equal(mocks.MockStorageNodes))
			gotBody, err := rlp.EncodeToBytes(converterPayload.BlockBody)
			Expect(err).ToNot(HaveOccurred())
			expectedBody, err := rlp.EncodeToBytes(mocks.MockBlock.Body())
			Expect(err).ToNot(HaveOccurred())
			Expect(gotBody).To(Equal(expectedBody))
			Expect(converterPayload.HeaderRLP).To(Equal(mocks.MockHeaderRlp))
			Expect(converterPayload.TrxMetaData).To(Equal(mocks.MockTrxMeta))
			Expect(converterPayload.ReceiptMetaData).To(Equal(mocks.MockRctMeta))
		})
		It(" Throws an error if the wrong chain config is used", func() {
			converter := ipfs.NewPayloadConverter(params.TestnetChainConfig)
			_, err := converter.Convert(mocks.MockStateDiffPayload)
			Expect(err).To(HaveOccurred())
		})
	})
})
