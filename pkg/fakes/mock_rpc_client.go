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

package fakes

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/p2p"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/core"
)

type MockRpcClient struct {
	callContextErr   error
	ipcPath          string
	nodeType         core.NodeType
	passedContext    context.Context
	passedMethod     string
	passedResult     interface{}
	returnPOAHeader  core.POAHeader
	supportedModules map[string]string
}

func NewMockRpcClient() *MockRpcClient {
	return &MockRpcClient{}
}

func (client *MockRpcClient) SetIpcPath(ipcPath string) {
	client.ipcPath = ipcPath
}

func (client *MockRpcClient) CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error {
	client.passedContext = ctx
	client.passedResult = result
	client.passedMethod = method
	switch method {
	case "admin_nodeInfo":
		if p, ok := result.(*p2p.NodeInfo); ok {
			p.ID = "enode://GethNode@172.17.0.1:30303"
			p.Name = "Geth/v1.7"
		}
	case "eth_getBlockByNumber":
		if p, ok := result.(*types.Header); ok {
			*p = types.Header{Number: big.NewInt(123)}
		}
		if p, ok := result.(*core.POAHeader); ok {
			*p = client.returnPOAHeader
		}
		if client.callContextErr != nil {
			return client.callContextErr
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
	return client.supportedModules, nil
}

func (client *MockRpcClient) SetSupporedModules(supportedModules map[string]string) {
	client.supportedModules = supportedModules
}

func (client *MockRpcClient) SetCallContextErr(err error) {
	client.callContextErr = err
}

func (client *MockRpcClient) SetReturnPOAHeader(header core.POAHeader) {
	client.returnPOAHeader = header
}

func (client *MockRpcClient) AssertCallContextCalledWith(ctx context.Context, result interface{}, method string) {
	Expect(client.passedContext).To(Equal(ctx))
	Expect(client.passedResult).To(BeAssignableToTypeOf(result))
	Expect(client.passedMethod).To(Equal(method))
}
