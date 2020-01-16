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
	"github.com/vulcanize/vulcanizedb/libraries/shared/streamer"
)

// IPLDResolver is the interface to resolving IPLDs
type IPLDResolver interface {
	ResolveIPLDs(ipfsBlocks IPLDWrapper) streamer.SuperNodePayload
	ResolveHeaders(iplds []blocks.Block) [][]byte
	ResolveUncles(iplds []blocks.Block) [][]byte
	ResolveTransactions(iplds []blocks.Block) [][]byte
	ResolveReceipts(blocks []blocks.Block) [][]byte
	ResolveState(iplds map[common.Hash]blocks.Block) map[common.Hash][]byte
	ResolveStorage(iplds map[common.Hash]map[common.Hash]blocks.Block) map[common.Hash]map[common.Hash][]byte
}

// EthIPLDResolver is the underlying struct to support the IPLDResolver interface
type EthIPLDResolver struct{}

// NewIPLDResolver returns a pointer to an EthIPLDResolver which satisfies the IPLDResolver interface
func NewIPLDResolver() *EthIPLDResolver {
	return &EthIPLDResolver{}
}

// ResolveIPLDs is the exported method for resolving all of the ETH IPLDs packaged in an IpfsBlockWrapper
func (eir *EthIPLDResolver) ResolveIPLDs(ipfsBlocks IPLDWrapper) streamer.SuperNodePayload {
	return streamer.SuperNodePayload{
		BlockNumber:     ipfsBlocks.BlockNumber,
		HeadersRlp:      eir.ResolveHeaders(ipfsBlocks.Headers),
		UnclesRlp:       eir.ResolveUncles(ipfsBlocks.Uncles),
		TransactionsRlp: eir.ResolveTransactions(ipfsBlocks.Transactions),
		ReceiptsRlp:     eir.ResolveReceipts(ipfsBlocks.Receipts),
		StateNodesRlp:   eir.ResolveState(ipfsBlocks.StateNodes),
		StorageNodesRlp: eir.ResolveStorage(ipfsBlocks.StorageNodes),
	}
}

func (eir *EthIPLDResolver) ResolveHeaders(iplds []blocks.Block) [][]byte {
	headerRlps := make([][]byte, 0, len(iplds))
	for _, ipld := range iplds {
		headerRlps = append(headerRlps, ipld.RawData())
	}
	return headerRlps
}

func (eir *EthIPLDResolver) ResolveUncles(iplds []blocks.Block) [][]byte {
	uncleRlps := make([][]byte, 0, len(iplds))
	for _, ipld := range iplds {
		uncleRlps = append(uncleRlps, ipld.RawData())
	}
	return uncleRlps
}

func (eir *EthIPLDResolver) ResolveTransactions(iplds []blocks.Block) [][]byte {
	trxs := make([][]byte, 0, len(iplds))
	for _, ipld := range iplds {
		trxs = append(trxs, ipld.RawData())
	}
	return trxs
}

func (eir *EthIPLDResolver) ResolveReceipts(iplds []blocks.Block) [][]byte {
	rcts := make([][]byte, 0, len(iplds))
	for _, ipld := range iplds {
		rcts = append(rcts, ipld.RawData())
	}
	return rcts
}

func (eir *EthIPLDResolver) ResolveState(iplds map[common.Hash]blocks.Block) map[common.Hash][]byte {
	stateNodes := make(map[common.Hash][]byte, len(iplds))
	for key, ipld := range iplds {
		stateNodes[key] = ipld.RawData()
	}
	return stateNodes
}

func (eir *EthIPLDResolver) ResolveStorage(iplds map[common.Hash]map[common.Hash]blocks.Block) map[common.Hash]map[common.Hash][]byte {
	storageNodes := make(map[common.Hash]map[common.Hash][]byte)
	for stateKey, storageIPLDs := range iplds {
		storageNodes[stateKey] = make(map[common.Hash][]byte)
		for storageKey, storageVal := range storageIPLDs {
			storageNodes[stateKey][storageKey] = storageVal.RawData()
		}
	}
	return storageNodes
}
