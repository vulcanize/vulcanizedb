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
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/seed_node"

	"github.com/vulcanize/vulcanizedb/pkg/ipfs"
	"github.com/vulcanize/vulcanizedb/pkg/ipfs/mocks"
)

var (
	resolver ipfs.IPLDResolver
)

var _ = Describe("Resolver", func() {
	Describe("ResolveIPLDs", func() {
		It("Resolves IPLD data to their correct geth data types and packages them to send to requesting transformers", func() {
			resolver = ipfs.NewIPLDResolver()
			seedNodePayload, err := resolver.ResolveIPLDs(mocks.MockIPLDWrapper)
			Expect(err).ToNot(HaveOccurred())
			Expect(seedNodePayload.BlockNumber.Int64()).To(Equal(int64(1)))
			Expect(seedNodePayload.HeadersRlp).To(Equal(mocks.MockSeeNodePayload.HeadersRlp))
			Expect(seedNodePayload.UnclesRlp).To(Equal(mocks.MockSeeNodePayload.UnclesRlp))
			Expect(len(seedNodePayload.TransactionsRlp)).To(Equal(2))
			Expect(seed_node.ListContainsBytes(seedNodePayload.TransactionsRlp, mocks.MockTransactions.GetRlp(0))).To(BeTrue())
			Expect(seed_node.ListContainsBytes(seedNodePayload.TransactionsRlp, mocks.MockTransactions.GetRlp(1))).To(BeTrue())
			Expect(len(seedNodePayload.ReceiptsRlp)).To(Equal(2))
			Expect(seed_node.ListContainsBytes(seedNodePayload.ReceiptsRlp, mocks.MockReceipts.GetRlp(0))).To(BeTrue())
			Expect(seed_node.ListContainsBytes(seedNodePayload.ReceiptsRlp, mocks.MockReceipts.GetRlp(1))).To(BeTrue())
			Expect(len(seedNodePayload.StateNodesRlp)).To(Equal(2))
			Expect(seedNodePayload.StorageNodesRlp).To(Equal(mocks.MockSeeNodePayload.StorageNodesRlp))
		})
	})
})
