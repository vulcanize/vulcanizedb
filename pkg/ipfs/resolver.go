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

package ipfs

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/ipfs/go-block-format"
)

// IPLDResolver is the interface to resolving IPLDs
type IPLDResolver interface {
	ResolveIPLDs(ipfsBlocks IpldWrapper) (*ResponsePayload, error)
}

// EthIPLDResolver is the underlying struct to support the IPLDResolver interface
type EthIPLDResolver struct{}

// NewIPLDResolver returns a pointer to an EthIPLDResolver which satisfies the IPLDResolver interface
func NewIPLDResolver() *EthIPLDResolver {
	return &EthIPLDResolver{}
}

// ResolveIPLDs is the exported method for resolving all of the ETH IPLDs packaged in an IpfsBlockWrapper
func (eir *EthIPLDResolver) ResolveIPLDs(ipfsBlocks IpldWrapper) (*ResponsePayload, error) {
	response := new(ResponsePayload)
	eir.resolveHeaders(ipfsBlocks.Headers, response)
	eir.resolveTransactions(ipfsBlocks.Transactions, response)
	eir.resolveReceipts(ipfsBlocks.Receipts, response)
	eir.resolveState(ipfsBlocks.StateNodes, response)
	eir.resolveStorage(ipfsBlocks.StorageNodes, response)
	return response, nil
}

func (eir *EthIPLDResolver) resolveHeaders(blocks []blocks.Block, response *ResponsePayload) {
	for _, block := range blocks {
		raw := block.RawData()
		response.HeadersRlp = append(response.HeadersRlp, raw)
	}
}

func (eir *EthIPLDResolver) resolveTransactions(blocks []blocks.Block, response *ResponsePayload) {
	for _, block := range blocks {
		raw := block.RawData()
		response.TransactionsRlp = append(response.TransactionsRlp, raw)
	}
}

func (eir *EthIPLDResolver) resolveReceipts(blocks []blocks.Block, response *ResponsePayload) {
	for _, block := range blocks {
		raw := block.RawData()
		response.ReceiptsRlp = append(response.ReceiptsRlp, raw)
	}
}

func (eir *EthIPLDResolver) resolveState(blocks map[common.Hash]blocks.Block, response *ResponsePayload) {
	if response.StateNodesRlp == nil {
		response.StateNodesRlp = make(map[common.Hash][]byte)
	}
	for key, block := range blocks {
		raw := block.RawData()
		response.StateNodesRlp[key] = raw
	}
}

func (eir *EthIPLDResolver) resolveStorage(blocks map[common.Hash]map[common.Hash]blocks.Block, response *ResponsePayload) {
	if response.StateNodesRlp == nil {
		response.StorageNodesRlp = make(map[common.Hash]map[common.Hash][]byte)
	}
	for stateKey, storageBlocks := range blocks {
		for storageKey, storageVal := range storageBlocks {
			raw := storageVal.RawData()
			response.StorageNodesRlp[stateKey][storageKey] = raw
		}
	}
}
