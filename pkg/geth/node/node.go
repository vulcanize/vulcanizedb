package node

import (
	"context"

	"strconv"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/rpc"
)

func Info(client *rpc.Client) core.Node {
	node := core.Node{}
	node.NetworkId = NetworkId(client)
	node.GenesisBlock = GenesisBlock(client)
	node.Id, node.ClientName = IdClientName(client)
	return node
}

func IdClientName(client *rpc.Client) (string, string) {
	var info p2p.NodeInfo
	modules, _ := client.SupportedModules()
	if _, ok := modules["admin"]; ok {
		client.CallContext(context.Background(), &info, "admin_nodeInfo")
		return info.ID, info.Name
	}
	return "", ""
}

func NetworkId(client *rpc.Client) float64 {
	var version string
	client.CallContext(context.Background(), &version, "net_version")
	networkId, _ := strconv.ParseFloat(version, 64)
	return networkId
}

func ProtocolVersion(client *rpc.Client) string {
	var protocolVersion string
	client.CallContext(context.Background(), &protocolVersion, "eth_protocolVersion")
	return protocolVersion
}

func GenesisBlock(client *rpc.Client) string {
	var header *types.Header
	client.CallContext(context.Background(), &header, "eth_getBlockByNumber", "0x0", false)
	return header.Hash().Hex()
}
