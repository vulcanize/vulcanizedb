// VulcanizeDB
// Copyright © 2019 Vulcanize

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.

// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package client

import (
	"context"
	"errors"
	"reflect"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/makerdao/vulcanizedb/pkg/core"
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

func (client RpcClient) BatchCall(batch []core.BatchElem) error {
	var rpcBatch []rpc.BatchElem
	for _, batchElem := range batch {
		var newBatchElem = rpc.BatchElem{
			Result: batchElem.Result,
			Method: batchElem.Method,
			Args:   batchElem.Args,
			Error:  batchElem.Error,
		}

		rpcBatch = append(rpcBatch, newBatchElem)
	}
	return client.client.BatchCall(rpcBatch)
}

// Subscribe subscribes to an rpc "namespace_subscribe" subscription with the given channel
// The first argument needs to be the method we wish to invoke
func (client RpcClient) Subscribe(namespace string, payloadChan interface{}, args ...interface{}) (core.Subscription, error) {
	chanVal := reflect.ValueOf(payloadChan)
	if chanVal.Kind() != reflect.Chan || chanVal.Type().ChanDir()&reflect.SendDir == 0 {
		return nil, errors.New("second argument to Subscribe must be a writable channel")
	}
	if chanVal.IsNil() {
		return nil, errors.New("channel given to Subscribe must not be nil")
	}
	rpcSubscription, err := client.client.Subscribe(context.Background(), namespace, payloadChan, args...)
	return Subscription{RpcSubscription: rpcSubscription}, err
}
