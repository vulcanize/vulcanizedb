package node

import (
	"context"

	"strconv"

	"regexp"

	"strings"

	"log"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/vulcanize/vulcanizedb/pkg/core"
)

type PropertiesReader interface {
	NodeInfo() (id string, name string)
	NetworkId() float64
	GenesisBlock() string
}

type ClientWrapper struct {
	ContextCaller
	IPCPath string
}

type ContextCaller interface {
	CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error
	SupportedModules() (map[string]string, error)
}

type ParityClient struct {
	ClientWrapper
}

type GethClient struct {
	ClientWrapper
}

type InfuraClient struct {
	ClientWrapper
}

func (clientWrapper ClientWrapper) NodeType() core.NodeType {
	if strings.Contains(clientWrapper.IPCPath, "infura") {
		return core.INFURA
	}
	modules, _ := clientWrapper.SupportedModules()
	if _, ok := modules["admin"]; ok {
		return core.GETH
	}
	return core.PARITY
}

func makePropertiesReader(wrapper ClientWrapper) PropertiesReader {
	switch wrapper.NodeType() {
	case core.GETH:
		return GethClient{ClientWrapper: wrapper}
	case core.PARITY:
		return ParityClient{ClientWrapper: wrapper}
	case core.INFURA:
		return InfuraClient{ClientWrapper: wrapper}
	default:
		return wrapper
	}
}

func MakeNode(wrapper ClientWrapper) core.Node {
	pr := makePropertiesReader(wrapper)
	id, name := pr.NodeInfo()
	return core.Node{
		GenesisBlock: pr.GenesisBlock(),
		NetworkID:    pr.NetworkId(),
		ID:           id,
		ClientName:   name,
	}
}

func (client ClientWrapper) NetworkId() float64 {
	var version string
	err := client.CallContext(context.Background(), &version, "net_version")
	if err != nil {
		log.Println(err)
	}
	networkId, _ := strconv.ParseFloat(version, 64)
	return networkId
}

func (client ClientWrapper) GenesisBlock() string {
	var header *types.Header
	blockZero := "0x0"
	includeTransactions := false
	client.CallContext(context.Background(), &header, "eth_getBlockByNumber", blockZero, includeTransactions)
	return header.Hash().Hex()
}

func (client ClientWrapper) NodeInfo() (string, string) {
	var info p2p.NodeInfo
	client.CallContext(context.Background(), &info, "admin_nodeInfo")
	return info.ID, info.Name
}

func (client ParityClient) NodeInfo() (string, string) {
	nodeInfo := client.parityNodeInfo()
	id := client.parityID()
	return id, nodeInfo
}

func (client InfuraClient) NodeInfo() (string, string) {
	return "infura", "infura"
}

func (client ParityClient) parityNodeInfo() string {
	var nodeInfo core.ParityNodeInfo
	client.CallContext(context.Background(), &nodeInfo, "parity_versionInfo")
	return nodeInfo.String()
}

func (client ParityClient) parityID() string {
	var enodeId = regexp.MustCompile(`^enode://(.+)@.+$`)
	var enodeURL string
	client.CallContext(context.Background(), &enodeURL, "parity_enode")
	enode := enodeId.FindStringSubmatch(enodeURL)
	if len(enode) < 2 {
		return ""
	}
	return enode[1]
}
