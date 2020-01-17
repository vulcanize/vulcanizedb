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

package eth

import (
	"context"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
)

// APIName is the namespace for the super node's eth api
const APIName = "eth"

// APIVersion is the version of the super node's eth api
const APIVersion = "0.0.1"

type PublicEthAPI struct {
	b *Backend
}

// NewPublicEthAPI creates a new PublicEthAPI with the provided underlying Backend
func NewPublicEthAPI(b *Backend) *PublicEthAPI {
	return &PublicEthAPI{
		b: b,
	}
}

/*
to start, need
eth_blockNumber
eth_getLogs
eth_getHeaderByNumber
*/

// BlockNumber returns the block number of the chain head.
func (pea *PublicEthAPI) BlockNumber() hexutil.Uint64 {
	number, _ := pea.b.retriever.RetrieveLastBlockNumber()
	return hexutil.Uint64(number)
}

// GetHeaderByNumber returns the requested canonical block header.
// * When blockNr is -1 the chain head is returned.
func (pea *PublicEthAPI) GetHeaderByNumber(ctx context.Context, number rpc.BlockNumber) (map[string]interface{}, error) {
	header, err := pea.b.HeaderByNumber(ctx, number)
	if header != nil && err == nil {
		return pea.rpcMarshalHeader(header), err
	}
	return nil, err
}

// rpcMarshalHeader uses the generalized output filler, then adds the total difficulty field, which requires
// a `PublicBlockchainAPI`.
func (pea *PublicEthAPI) rpcMarshalHeader(header *types.Header) map[string]interface{} {
	fields := RPCMarshalHeader(header)
	fields["totalDifficulty"] = (*hexutil.Big)(pea.b.GetTd(header.Hash()))
	return fields
}

// RPCMarshalHeader converts the given header to the RPC output .
func RPCMarshalHeader(head *types.Header) map[string]interface{} {
	return map[string]interface{}{
		"number":           (*hexutil.Big)(head.Number),
		"hash":             head.Hash(),
		"parentHash":       head.ParentHash,
		"nonce":            head.Nonce,
		"mixHash":          head.MixDigest,
		"sha3Uncles":       head.UncleHash,
		"logsBloom":        head.Bloom,
		"stateRoot":        head.Root,
		"miner":            head.Coinbase,
		"difficulty":       (*hexutil.Big)(head.Difficulty),
		"extraData":        hexutil.Bytes(head.Extra),
		"size":             hexutil.Uint64(head.Size()),
		"gasLimit":         hexutil.Uint64(head.GasLimit),
		"gasUsed":          hexutil.Uint64(head.GasUsed),
		"timestamp":        hexutil.Uint64(head.Time),
		"transactionsRoot": head.TxHash,
		"receiptsRoot":     head.ReceiptHash,
	}
}
