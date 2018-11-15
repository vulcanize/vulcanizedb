// VulcanizeDB
// Copyright © 2018 Vulcanize

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
