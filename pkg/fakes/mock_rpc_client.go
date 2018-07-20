package fakes

import (
	"context"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/vulcanize/vulcanizedb/pkg/core"
)

type MockRpcClient struct {
	ipcPath  string
	nodeType core.NodeType
}

func NewMockRpcClient() *MockRpcClient {
	return &MockRpcClient{}
}

func (client *MockRpcClient) SetNodeType(nodeType core.NodeType) {
	client.nodeType = nodeType
}

func (client *MockRpcClient) SetIpcPath(ipcPath string) {
	client.ipcPath = ipcPath
}

func (*MockRpcClient) CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error {
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

func (client *MockRpcClient) IpcPath() string {
	return client.ipcPath
}

func (client *MockRpcClient) SupportedModules() (map[string]string, error) {
	result := make(map[string]string)
	if client.nodeType == core.GETH {
		result["admin"] = "ok"
	}
	return result, nil
}
