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

package fakes

import (
	"context"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/ethereum/go-ethereum/statediff"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/eth/client"
	. "github.com/onsi/gomega"
)

type MockRpcClient struct {
	callContextErr      error
	ClientVersion       string
	GethNodeInfo        p2p.NodeInfo
	ipcPath             string
	NetworkID           string
	nodeType            core.NodeType
	ParityEnode         string
	ParityNodeInfo      core.ParityNodeInfo
	passedContext       context.Context
	passedMethod        string
	passedResult        interface{}
	passedBatch         []core.BatchElem
	passedNamespace     string
	passedPayloadChan   chan statediff.Payload
	passedSubscribeArgs []interface{}
	lengthOfBatch       int
	returnPOAHeader     core.POAHeader
	returnPOAHeaders    []core.POAHeader
	returnPOWHeaders    []*types.Header
}

func NewMockRpcClient() *MockRpcClient {
	return &MockRpcClient{}
}

func (c *MockRpcClient) Subscribe(namespace string, payloadChan interface{}, args ...interface{}) (core.Subscription, error) {
	c.passedNamespace = namespace

	passedPayloadChan, ok := payloadChan.(chan statediff.Payload)
	if !ok {
		return nil, errors.New("passed in channel is not of the correct type")
	}
	c.passedPayloadChan = passedPayloadChan

	for _, arg := range args {
		c.passedSubscribeArgs = append(c.passedSubscribeArgs, arg)
	}

	subscription := rpc.ClientSubscription{}
	return client.Subscription{RpcSubscription: &subscription}, nil
}

func (c *MockRpcClient) AssertSubscribeCalledWith(namespace string, payloadChan chan statediff.Payload, args []interface{}) {
	Expect(c.passedNamespace).To(Equal(namespace))
	Expect(c.passedPayloadChan).To(Equal(payloadChan))
	Expect(c.passedSubscribeArgs).To(Equal(args))
}

func (c *MockRpcClient) SetIpcPath(ipcPath string) {
	c.ipcPath = ipcPath
}

func (c *MockRpcClient) BatchCall(batch []core.BatchElem) error {
	c.passedBatch = batch
	c.passedMethod = batch[0].Method
	c.lengthOfBatch = len(batch)

	for _, batchElem := range batch {
		c.passedContext = context.Background()
		c.passedResult = &batchElem.Result
		c.passedMethod = batchElem.Method
		if p, ok := batchElem.Result.(*types.Header); ok {
			*p = types.Header{Number: big.NewInt(100)}
		}
		if p, ok := batchElem.Result.(*core.POAHeader); ok {
			*p = c.returnPOAHeader
		}
	}

	return nil
}

func (c *MockRpcClient) CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error {
	c.passedContext = ctx
	c.passedResult = result
	c.passedMethod = method
	switch method {
	case "eth_getBlockByNumber":
		if p, ok := result.(*types.Header); ok {
			*p = types.Header{Number: big.NewInt(100)}
		}
		if p, ok := result.(*core.POAHeader); ok {

			*p = c.returnPOAHeader
		}
		if c.callContextErr != nil {
			return c.callContextErr
		}
	case "parity_versionInfo":
		if p, ok := result.(*core.ParityNodeInfo); ok {
			*p = c.ParityNodeInfo
		}
	case "parity_enode":
		if p, ok := result.(*string); ok {
			*p = c.ParityEnode
		}
	case "net_version":
		if p, ok := result.(*string); ok {
			*p = c.NetworkID
		}
	case "web3_clientVersion":
		if p, ok := result.(*string); ok {
			*p = c.ClientVersion
		}
	}
	return nil
}

func (c *MockRpcClient) IpcPath() string {
	return c.ipcPath
}

func (c *MockRpcClient) SetCallContextErr(err error) {
	c.callContextErr = err
}

func (c *MockRpcClient) SetReturnPOAHeader(header core.POAHeader) {
	c.returnPOAHeader = header
}

func (c *MockRpcClient) SetReturnPOWHeaders(headers []*types.Header) {
	c.returnPOWHeaders = headers
}

func (c *MockRpcClient) SetReturnPOAHeaders(headers []core.POAHeader) {
	c.returnPOAHeaders = headers
}

func (c *MockRpcClient) AssertCallContextCalledWith(ctx context.Context, result interface{}, method string) {
	Expect(c.passedContext).To(Equal(ctx))
	Expect(c.passedResult).To(BeAssignableToTypeOf(result))
	Expect(c.passedMethod).To(Equal(method))
}

func (c *MockRpcClient) AssertBatchCalledWith(method string, lengthOfBatch int) {
	Expect(c.lengthOfBatch).To(Equal(lengthOfBatch))
	for _, batch := range c.passedBatch {
		Expect(batch.Method).To(Equal(method))
	}
	Expect(c.passedMethod).To(Equal(method))
}
