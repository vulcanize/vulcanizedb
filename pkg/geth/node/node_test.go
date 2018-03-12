package node_test

import (
	"encoding/json"

	"context"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/p2p"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/geth/node"
)

type MockContextCaller struct {
	nodeType core.NodeType
}

var EmpytHeaderHash = "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347"

func (MockContextCaller) CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error {
	switch method {
	case "admin_nodeInfo":
		if p, ok := result.(*p2p.NodeInfo); ok {
			p.ID = "enode://GethNode@172.17.0.1:30303"
			p.Name = "Geth/v1.7"
		}
	case "eth_getBlockByNumber":
		if p, ok := result.(*types.Header); ok {
			*p = types.Header{}
		}

	case "parity_versionInfo":
		if p, ok := result.(*core.ParityNodeInfo); ok {
			*p = core.ParityNodeInfo{
				Track: "",
				ParityVersion: core.ParityVersion{
					Major: 1,
					Minor: 2,
					Patch: 3,
				},
				Hash: "",
			}
		}
	case "parity_enode":
		if p, ok := result.(*string); ok {
			*p = "enode://ParityNode@172.17.0.1:30303"
		}
	case "net_version":
		if p, ok := result.(*string); ok {
			*p = "1234"
		}
	}
	return nil
}

func (mcc MockContextCaller) SupportedModules() (map[string]string, error) {
	result := make(map[string]string)
	if mcc.nodeType == core.GETH {
		result["admin"] = "ok"
	}
	return result, nil
}

var _ = Describe("Parity Node Info", func() {

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

	It("returns the genesis block for any client", func() {
		mcc := MockContextCaller{}
		cw := node.ClientWrapper{ContextCaller: mcc}
		n := node.MakeNode(cw)
		Expect(n.GenesisBlock).To(Equal(EmpytHeaderHash))
	})

	It("returns the network id for any client", func() {
		mcc := MockContextCaller{}
		cw := node.ClientWrapper{ContextCaller: mcc}
		n := node.MakeNode(cw)
		Expect(n.NetworkID).To(Equal(float64(1234)))
	})

	It("returns parity ID and client name for parity node", func() {
		mcc := MockContextCaller{core.PARITY}
		cw := node.ClientWrapper{ContextCaller: mcc}
		n := node.MakeNode(cw)
		Expect(n.ID).To(Equal("ParityNode"))
		Expect(n.ClientName).To(Equal("Parity/v1.2.3/"))
	})

	It("returns geth ID and client name for geth node", func() {
		mcc := MockContextCaller{core.GETH}
		cw := node.ClientWrapper{ContextCaller: mcc}
		n := node.MakeNode(cw)
		Expect(n.ID).To(Equal("enode://GethNode@172.17.0.1:30303"))
		Expect(n.ClientName).To(Equal("Geth/v1.7"))
	})

	It("returns infura ID and client name for infura node", func() {
		mcc := MockContextCaller{core.INFURA}
		cw := node.ClientWrapper{ContextCaller: mcc, IPCPath: "https://mainnet.infura.io/123"}
		n := node.MakeNode(cw)
		Expect(n.ID).To(Equal("infura"))
		Expect(n.ClientName).To(Equal("infura"))
	})
})
