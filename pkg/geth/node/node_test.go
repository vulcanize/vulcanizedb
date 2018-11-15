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

package node_test

import (
	"encoding/json"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/geth/node"
)

var EmpytHeaderHash = "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347"

var _ = Describe("Node Info", func() {
	Describe("Parity Node Info", func() {
		It("verifies parity_versionInfo can be unmarshalled into ParityNodeInfo", func() {
			var parityNodeInfo core.ParityNodeInfo
			nodeInfoJSON := []byte(
				`{
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

		It("returns parity ID and client name for parity node", func() {
			client := fakes.NewMockRpcClient()

			n := node.MakeNode(client)
			Expect(n.ID).To(Equal("ParityNode"))
			Expect(n.ClientName).To(Equal("Parity/v1.2.3/"))
		})
	})

	It("returns the genesis block for any client", func() {
		client := fakes.NewMockRpcClient()
		n := node.MakeNode(client)
		Expect(n.GenesisBlock).To(Equal(EmpytHeaderHash))
	})

	It("returns the network id for any client", func() {
		client := fakes.NewMockRpcClient()
		n := node.MakeNode(client)
		Expect(n.NetworkID).To(Equal(float64(1234)))
	})

	It("returns geth ID and client name for geth node", func() {
		client := fakes.NewMockRpcClient()
		supportedModules := make(map[string]string)
		supportedModules["admin"] = "ok"
		client.SetSupporedModules(supportedModules)

		n := node.MakeNode(client)
		Expect(n.ID).To(Equal("enode://GethNode@172.17.0.1:30303"))
		Expect(n.ClientName).To(Equal("Geth/v1.7"))
	})

	It("returns infura ID and client name for infura node", func() {
		client := fakes.NewMockRpcClient()
		client.SetIpcPath("infura/path")
		n := node.MakeNode(client)
		Expect(n.ID).To(Equal("infura"))
		Expect(n.ClientName).To(Equal("infura"))
	})

	It("returns local id and client name for Local node", func() {
		client := fakes.NewMockRpcClient()
		client.SetIpcPath("127.0.0.1")
		n := node.MakeNode(client)
		Expect(n.ID).To(Equal("ganache"))
		Expect(n.ClientName).To(Equal("ganache"))

		client.SetIpcPath("localhost")
		n = node.MakeNode(client)
		Expect(n.ID).To(Equal("ganache"))
		Expect(n.ClientName).To(Equal("ganache"))
	})
})
