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

package node_test

import (
	"encoding/json"
	"strconv"

	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/eth/node"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var EmptyHeaderHash = "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347"

var _ = Describe("Node Info", func() {
	Describe("Parity Node Info", func() {
		It("verifies parity_versionInfo can be unmarshalled into ParityNodeInfo", func() {
			var parityNodeInfo core.ParityNodeInfo
			nodeInfoJSON := []byte(`{
				"hash": "0x2ae8b4ca278dd7b896090366615fef81cbbbc0e0",
				"track": "null",
				"version": {
				"major": 1,
				"minor": 6,
				"patch": 0
				}
			}`)
			json.Unmarshal(nodeInfoJSON, &parityNodeInfo)
			Expect(parityNodeInfo.Hash).To(Equal("0x2ae8b4ca278dd7b896090366615fef81cbbbc0e0"))
			Expect(parityNodeInfo.Track).To(Equal("null"))
			Expect(parityNodeInfo.Major).To(Equal(1))
			Expect(parityNodeInfo.Minor).To(Equal(6))
			Expect(parityNodeInfo.Patch).To(Equal(0))
		})

		It("Creates client string", func() {
			parityNodeInfo := core.ParityNodeInfo{
				Track: "null",
				ParityVersion: core.ParityVersion{
					Major: 1,
					Minor: 6,
					Patch: 0,
				},
				Hash: "0x1232144j",
			}
			Expect(parityNodeInfo.String()).To(Equal("Parity/v1.6.0/"))
		})

		It("returns client name for parity node", func() {
			client := fakes.NewMockRpcClient()
			client.ClientVersion = "Parity-Ethereum//v2.5.13-stable-253ff3f-20191231/x86_64-linux-gnu/rustc1.40.0"
			client.ParityEnode = "enode://ParityNode@172.17.0.1:30303"
			client.ParityNodeInfo = core.ParityNodeInfo{
				ParityVersion: core.ParityVersion{
					Major: 1,
					Minor: 2,
					Patch: 3,
				},
			}

			n := node.MakeNode(client)

			Expect(n.ClientName).To(Equal("Parity/v1.2.3/"))
		})
	})

	It("returns the genesis block for any client", func() {
		client := fakes.NewMockRpcClient()

		n := node.MakeNode(client)

		Expect(n.GenesisBlock).To(Equal(EmptyHeaderHash))
	})

	It("returns the network id for any client", func() {
		client := fakes.NewMockRpcClient()
		client.NetworkID = "1234"

		n := node.MakeNode(client)

		expectedNetworkID, err := strconv.ParseFloat(client.NetworkID, 64)
		Expect(err).NotTo(HaveOccurred())
		Expect(n.NetworkID).To(Equal(expectedNetworkID))
	})

	It("returns the IpcPath as node ID", func() {
		client := fakes.NewMockRpcClient()
		client.SetIpcPath("infura/path")

		n := node.MakeNode(client)

		Expect(n.ID).To(Equal(client.IpcPath()))
	})

	It("returns client name for geth node", func() {
		client := fakes.NewMockRpcClient()
		client.ClientVersion = "Geth/v1.9.9-omnibus-e320ae4c-20191206/linux-amd64/go1.13.4"

		n := node.MakeNode(client)

		Expect(n.ClientName).To(Equal(client.ClientVersion))
	})

	It("returns client name for infura node", func() {
		client := fakes.NewMockRpcClient()
		client.SetIpcPath("infura/path")

		n := node.MakeNode(client)

		Expect(n.ClientName).To(Equal("infura"))
	})

	It("returns ganache by default", func() {
		client := fakes.NewMockRpcClient()

		n := node.MakeNode(client)

		Expect(n.ClientName).To(Equal("ganache"))
	})
})
