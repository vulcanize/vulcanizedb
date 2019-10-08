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

	"github.com/vulcanize/vulcanizedb/pkg/ipfs"
	"github.com/vulcanize/vulcanizedb/pkg/ipfs/mocks"
	"github.com/vulcanize/vulcanizedb/pkg/super_node"
)

var (
	resolver ipfs.IPLDResolver
)

var _ = Describe("Resolver", func() {
	Describe("ResolveIPLDs", func() {
		BeforeEach(func() {
			resolver = ipfs.NewIPLDResolver()
		})
		It("Resolves IPLD data to their correct geth data types and packages them to send to requesting transformers", func() {
			superNodePayload := resolver.ResolveIPLDs(mocks.MockIPLDWrapper)
			Expect(superNodePayload.BlockNumber.Int64()).To(Equal(mocks.MockSeeNodePayload.BlockNumber.Int64()))
			Expect(superNodePayload.HeadersRlp).To(Equal(mocks.MockSeeNodePayload.HeadersRlp))
			Expect(superNodePayload.UnclesRlp).To(Equal(mocks.MockSeeNodePayload.UnclesRlp))
			Expect(len(superNodePayload.TransactionsRlp)).To(Equal(2))
			Expect(super_node.ListContainsBytes(superNodePayload.TransactionsRlp, mocks.MockTransactions.GetRlp(0))).To(BeTrue())
			Expect(super_node.ListContainsBytes(superNodePayload.TransactionsRlp, mocks.MockTransactions.GetRlp(1))).To(BeTrue())
			Expect(len(superNodePayload.ReceiptsRlp)).To(Equal(2))
			Expect(super_node.ListContainsBytes(superNodePayload.ReceiptsRlp, mocks.MockReceipts.GetRlp(0))).To(BeTrue())
			Expect(super_node.ListContainsBytes(superNodePayload.ReceiptsRlp, mocks.MockReceipts.GetRlp(1))).To(BeTrue())
			Expect(len(superNodePayload.StateNodesRlp)).To(Equal(2))
			Expect(superNodePayload.StorageNodesRlp).To(Equal(mocks.MockSeeNodePayload.StorageNodesRlp))
		})
	})
})
