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

package btc

import (
	"encoding/json"
	"math/big"

	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/ipfs/go-block-format"
)

// BlockPayload packages the block and tx data received from block connection notifications
type BlockPayload struct {
	Height int64
	Header *wire.BlockHeader
	Txs    []*btcutil.Tx
}

// IPLDPayload is a custom type which packages raw BTC data for publishing to IPFS and filtering to subscribers
// Returned by PayloadConverter
// Passed to IPLDPublisher and ResponseFilterer
type IPLDPayload struct {
	BlockPayload
	TxMetaData []TxModelWithInsAndOuts
}

// CIDPayload is a struct to hold all the CIDs and their associated meta data for indexing in Postgres
// Returned by IPLDPublisher
// Passed to CIDIndexer
type CIDPayload struct {
	HeaderCID       HeaderModel
	TransactionCIDs []TxModelWithInsAndOuts
}

// CIDWrapper is used to direct fetching of IPLDs from IPFS
// Returned by CIDRetriever
// Passed to IPLDFetcher
type CIDWrapper struct {
	BlockNumber  *big.Int
	Headers      []HeaderModel
	Transactions []TxModel
}

// IPLDWrapper is used to package raw IPLD block data fetched from IPFS
// Returned by IPLDFetcher
// Passed to IPLDResolver
type IPLDWrapper struct {
	BlockNumber  *big.Int
	Headers      []blocks.Block
	Transactions []blocks.Block
}

// StreamResponse holds the data streamed from the super node eth service to the requesting clients
// Returned by IPLDResolver and ResponseFilterer
// Passed to client subscriptions
type StreamResponse struct {
	BlockNumber       *big.Int `json:"blockNumber"`
	SerializedHeaders [][]byte `json:"headerBytes"`
	SerializedTxs     [][]byte `json:"transactionBytes"`

	encoded []byte
	err     error
}

func (sr *StreamResponse) ensureEncoded() {
	if sr.encoded == nil && sr.err == nil {
		sr.encoded, sr.err = json.Marshal(sr)
	}
}

// Length to implement Encoder interface for StateDiff
func (sr *StreamResponse) Length() int {
	sr.ensureEncoded()
	return len(sr.encoded)
}

// Encode to implement Encoder interface for StateDiff
func (sr *StreamResponse) Encode() ([]byte, error) {
	sr.ensureEncoded()
	return sr.encoded, sr.err
}
