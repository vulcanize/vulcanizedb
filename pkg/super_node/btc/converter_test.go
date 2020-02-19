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

package btc_test

import (
	"github.com/btcsuite/btcd/chaincfg"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/super_node/btc"
	"github.com/vulcanize/vulcanizedb/pkg/super_node/btc/mocks"
)

var _ = Describe("Converter", func() {
	Describe("Convert", func() {
		It("Converts mock BlockPayloads into the expected IPLDPayloads", func() {
			converter := btc.NewPayloadConverter(&chaincfg.MainNetParams)
			payload, err := converter.Convert(mocks.MockBlockPayload)
			Expect(err).ToNot(HaveOccurred())
			convertedPayload, ok := payload.(btc.IPLDPayload)
			Expect(ok).To(BeTrue())
			Expect(convertedPayload).To(Equal(mocks.MockIPLDPayload))
			Expect(convertedPayload.BlockHeight).To(Equal(mocks.MockBlockHeight))
			Expect(convertedPayload.Header).To(Equal(&mocks.MockBlock.Header))
			Expect(convertedPayload.Txs).To(Equal(mocks.MockTransactions))
			Expect(convertedPayload.TxMetaData).To(Equal(mocks.MockTxsMetaData))
		})
	})
})
