package node

import (
	"context"

	"strconv"

	"github.com/8thlight/vulcanizedb/pkg/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
)

func Retrieve(client *rpc.Client) core.Node {
	node := core.Node{}

	var version string
	client.CallContext(context.Background(), &version, "net_version")
	node.NetworkId, _ = strconv.ParseFloat(version, 64)

	var protocolVersion string
	client.CallContext(context.Background(), &protocolVersion, "eth_protocolVersion")

	var header *types.Header
	client.CallContext(context.Background(), &header, "eth_getBlockByNumber", "0x0", false)
	node.GenesisBlock = header.Hash().Hex()

	return node
}
