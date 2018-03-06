package node

import (
	"context"

	"strconv"

	"regexp"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/vulcanize/vulcanizedb/pkg/core"
)

type Getter interface {
	NetworkId() float64
	GenesisBlock() string
	GethNodeInfo() (string, string)
	ParityNodeInfo() (string, string)
}

type NodeFactory struct {
	Client *rpc.Client
}

func MakeNode(getter Getter) core.Node {
	node := core.Node{}
	node.NetworkID = getter.NetworkId()
	node.GenesisBlock = getter.GenesisBlock()
	node.ID, node.ClientName = getter.GethNodeInfo()
	if notGethNode(node) {
		node.ID, node.ClientName = getter.ParityNodeInfo()
	}
	return node
}

func notGethNode(node core.Node) bool {
	return node.ID == ""
}

func (nf *NodeFactory) GethNodeInfo() (string, string) {
	var info p2p.NodeInfo
	modules, _ := nf.Client.SupportedModules()
	if _, ok := modules["admin"]; ok {
		nf.Client.CallContext(context.Background(), &info, "admin_nodeInfo")
		return info.ID, info.Name
	}
	return "", ""
}

func (nf *NodeFactory) NetworkId() float64 {
	var version string
	nf.Client.CallContext(context.Background(), &version, "net_version")
	networkId, _ := strconv.ParseFloat(version, 64)
	return networkId
}

func (nf *NodeFactory) ParityNodeInfo() (string, string) {
	client := nf.parityClient()
	id := nf.parityId()
	return id, client
}

func (nf *NodeFactory) parityClient() string {
	var nodeInfo core.ParityNodeInfo
	nf.Client.CallContext(context.Background(), &nodeInfo, "parity_versionInfo")
	return nodeInfo.String()
}

func (nf *NodeFactory) parityId() string {
	var enodeId = regexp.MustCompile(`^enode://(.+)@.+$`)
	var enodeURL string
	nf.Client.CallContext(context.Background(), &enodeURL, "parity_enode")
	enode := enodeId.FindStringSubmatch(enodeURL)
	if len(enode) < 2 {
		return ""
	}
	return enode[1]
}

func (nf *NodeFactory) GenesisBlock() string {
	var header *types.Header
	nf.Client.CallContext(context.Background(), &header, "eth_getBlockByNumber", "0x0", false)
	return header.Hash().Hex()
}
