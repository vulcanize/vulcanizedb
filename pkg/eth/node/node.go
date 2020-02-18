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

package node

import (
	"context"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/sirupsen/logrus"
)

type IPropertiesReader interface {
	NodeClientName() string
	NetworkId() float64
	GenesisBlock() string
}

type PropertiesReader struct {
	client core.RpcClient
}

type ParityClient struct {
	PropertiesReader
}

type GethClient struct {
	PropertiesReader
}

type InfuraClient struct {
	PropertiesReader
}

type GanacheClient struct {
	PropertiesReader
}

func MakeNode(rpcClient core.RpcClient) core.Node {
	pr := makePropertiesReader(rpcClient)
	return core.Node{
		GenesisBlock: pr.GenesisBlock(),
		NetworkID:    pr.NetworkId(),
		ID:           rpcClient.IpcPath(),
		ClientName:   pr.NodeClientName(),
	}
}

func makePropertiesReader(client core.RpcClient) IPropertiesReader {
	switch getNodeType(client) {
	case core.GETH:
		return GethClient{PropertiesReader: PropertiesReader{client: client}}
	case core.PARITY:
		return ParityClient{PropertiesReader: PropertiesReader{client: client}}
	case core.INFURA:
		return InfuraClient{PropertiesReader: PropertiesReader{client: client}}
	case core.GANACHE:
		return GanacheClient{PropertiesReader: PropertiesReader{client: client}}
	default:
		return PropertiesReader{client: client}
	}
}

func getNodeType(client core.RpcClient) core.NodeType {
	if strings.Contains(client.IpcPath(), "infura") {
		return core.INFURA
	}
	var version string
	err := client.CallContext(context.Background(), &version, "web3_clientVersion")
	if err != nil {
		logrus.Warnf("error getting client version: %s\n", err.Error())
	}
	if strings.Contains(version, "Geth") {
		return core.GETH
	}
	if strings.Contains(version, "Parity") {
		return core.PARITY
	}
	return core.GANACHE
}

func (reader PropertiesReader) NetworkId() float64 {
	var version string
	err := reader.client.CallContext(context.Background(), &version, "net_version")
	if err != nil {
		logrus.Warnf("error getting net_version: %s", err.Error())
	}
	networkId, _ := strconv.ParseFloat(version, 64)
	return networkId
}

func (reader PropertiesReader) GenesisBlock() string {
	var header *types.Header
	blockZero := "0x0"
	includeTransactions := false
	err := reader.client.CallContext(context.Background(), &header, "eth_getBlockByNumber", blockZero, includeTransactions)
	if err != nil {
		logrus.Warnf("error getting genesis block: %s", err.Error())
	}
	return header.Hash().Hex()
}

func (reader PropertiesReader) NodeClientName() string {
	var nodeName string
	err := reader.client.CallContext(context.Background(), &nodeName, "web3_clientVersion")
	if err != nil {
		logrus.Debugf("error getting client version: %s", err.Error())
	}
	return nodeName
}

func (client ParityClient) NodeClientName() string {
	return client.parityNodeInfo()
}

func (client InfuraClient) NodeClientName() string {
	return "infura"
}

func (client GanacheClient) NodeClientName() string {
	return "ganache"
}

func (client ParityClient) parityNodeInfo() string {
	var nodeInfo core.ParityNodeInfo
	err := client.client.CallContext(context.Background(), &nodeInfo, "parity_versionInfo")
	if err != nil {
		logrus.Warnf("error getting parity_versionInfo: %s", err.Error())
	}
	return nodeInfo.String()
}
