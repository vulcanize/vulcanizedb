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

	"github.com/vulcanize/vulcanizedb/pkg/super_node/shared"

	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/ipfs/go-block-format"
)

// BlockPayload packages the block and tx data received from block connection notifications
type BlockPayload struct {
	Height int32
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

func (ip IPLDPayload) Value() shared.StreamedIPLDs {
	return ip
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

// StreamPayload holds the data streamed from the super node eth service to the requesting clients
// Returned by IPLDResolver and ResponseFilterer
// Passed to client subscriptions
type StreamPayload struct {
	BlockNumber       *big.Int `json:"blockNumber"`
	HeadersBytes      [][]byte `json:"headerBytes"`
	TransactionsBytes [][]byte `json:"transactionBytes"`

	encoded []byte
	err     error
}

func (sd *StreamPayload) ensureEncoded() {
	if sd.encoded == nil && sd.err == nil {
		sd.encoded, sd.err = json.Marshal(sd)
	}
}

// Length to implement Encoder interface for StateDiff
func (sd *StreamPayload) Length() int {
	sd.ensureEncoded()
	return len(sd.encoded)
}

// Encode to implement Encoder interface for StateDiff
func (sd *StreamPayload) Encode() ([]byte, error) {
	sd.ensureEncoded()
	return sd.encoded, sd.err
}
