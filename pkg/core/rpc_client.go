package core

import "context"

type RpcClient interface {
	CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error
	IpcPath() string
	SupportedModules() (map[string]string, error)
}
