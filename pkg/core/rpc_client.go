package core

import (
	"context"
	"github.com/ethereum/go-ethereum/rpc"
)

type RpcClient interface {
	CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error
	BatchCall(batch []rpc.BatchElem) error
	IpcPath() string
	SupportedModules() (map[string]string, error)
}
