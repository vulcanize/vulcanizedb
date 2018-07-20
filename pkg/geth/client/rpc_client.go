package client

import (
	"context"
	"github.com/ethereum/go-ethereum/rpc"
)

type RpcClient struct {
	client  *rpc.Client
	ipcPath string
}

func NewRpcClient(client *rpc.Client, ipcPath string) RpcClient {
	return RpcClient{
		client:  client,
		ipcPath: ipcPath,
	}
}

func (client RpcClient) CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error {
	return client.client.CallContext(ctx, result, method, args)
}

func (client RpcClient) IpcPath() string {
	return client.ipcPath
}

func (client RpcClient) SupportedModules() (map[string]string, error) {
	return client.client.SupportedModules()
}
