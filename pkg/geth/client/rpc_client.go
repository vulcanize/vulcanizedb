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
	//If an empty interface (or other nil object) is passed to CallContext, when the JSONRPC message is created the params will
	//be interpreted as [null]. This seems to work fine for most of the ethereum clients (which presumably ignore a null parameter.
	//Ganache however does not ignore it, and throws an 'Incorrect number of arguments' error.
	if args == nil {
		return client.client.CallContext(ctx, result, method)
	} else {
		return client.client.CallContext(ctx, result, method, args...)
	}
}

func (client RpcClient) IpcPath() string {
	return client.ipcPath
}

func (client RpcClient) SupportedModules() (map[string]string, error) {
	return client.client.SupportedModules()
}

func (client RpcClient) BatchCall(batch []rpc.BatchElem) error {
	return client.client.BatchCall(batch)
}
