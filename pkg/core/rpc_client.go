package core

import (
	"context"
	"github.com/vulcanize/vulcanizedb/pkg/geth/client"
)

type RpcClient interface {
	CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error
	BatchCall(batch []client.BatchElem) error
	IpcPath() string
	SupportedModules() (map[string]string, error)
}
