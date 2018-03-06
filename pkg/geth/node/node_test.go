package node_test

import (
	"encoding/json"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/geth/node"
)

type MockNodeFactory struct {
	gethNodeName   string
	gethNodeId     string
	parityNodeName string
	parityNodeId   string
	networkId      float64
	genesisBlock   string
}

func (mnf MockNodeFactory) NetworkId() float64 {
	return mnf.networkId
}

func (mnf MockNodeFactory) GenesisBlock() string {
	return mnf.genesisBlock
}

func (mnf MockNodeFactory) GethNodeInfo() (string, string) {
	return mnf.gethNodeId, mnf.gethNodeName
}

func (mnf MockNodeFactory) ParityNodeInfo() (string, string) {
	return mnf.parityNodeId, mnf.parityNodeName
}

var _ = Describe("Parity Node Info", func() {

	It("Decodes json from parity_versionInfo", func() {
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

	It("Returns parity node for parity client", func() {
		mf := MockNodeFactory{parityNodeId: "0x1232389", parityNodeName: "Parity"}
		node := node.MakeNode(mf)
		Expect(node.ClientName).To(Equal("Parity"))
		Expect(node.ID).To(Equal("0x1232389"))
	})

	It("Returns geth node for geth client", func() {
		mf := MockNodeFactory{gethNodeId: "0x1234", gethNodeName: "Geth"}
		node := node.MakeNode(mf)
		Expect(node.ClientName).To(Equal("Geth"))
		Expect(node.ID).To(Equal("0x1234"))
	})

})
