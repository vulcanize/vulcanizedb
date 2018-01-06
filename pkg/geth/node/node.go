package node

import (
	"context"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/rpc"
)

func Retrieve(client *rpc.Client) core.Node {
	var info p2p.NodeInfo
	node := core.Node{}
	client.CallContext(context.Background(), &info, "admin_nodeInfo")
	for protocolName, protocol := range info.Protocols {
		if protocolName == "eth" {
			protocolMap, _ := protocol.(map[string]interface{})
			node.GenesisBlock = getAttribute(protocolMap, "genesis").(string)
			node.NetworkId = getAttribute(protocolMap, "network").(float64)
		}
	}
	return node
}

func getAttribute(protocolMap map[string]interface{}, protocol string) interface{} {
	for key, val := range protocolMap {
		if key == protocol {
			return val
		}
	}
	return nil
}
