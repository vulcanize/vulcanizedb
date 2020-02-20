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
	"math/big"

	"github.com/vulcanize/vulcanizedb/pkg/ipfs"

	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
)

// BlockPayload packages the block and tx data received from block connection notifications
type BlockPayload struct {
	BlockHeight int64
	Header      *wire.BlockHeader
	Txs         []*btcutil.Tx
}

// ConvertedPayload is a custom type which packages raw BTC data for publishing to IPFS and filtering to subscribers
// Returned by PayloadConverter
// Passed to IPLDPublisher and ResponseFilterer
type ConvertedPayload struct {
	BlockPayload
	TxMetaData []TxModelWithInsAndOuts
}

// Height satisfies the StreamedIPLDs interface
func (i ConvertedPayload) Height() int64 {
	return i.BlockPayload.BlockHeight
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

// IPLDs is used to package raw IPLD block data fetched from IPFS and returned by the server
// Returned by IPLDFetcher and ResponseFilterer
type IPLDs struct {
	BlockNumber  *big.Int
	Headers      []ipfs.BlockModel
	Transactions []ipfs.BlockModel
}

// Height satisfies the StreamedIPLDs interface
func (i IPLDs) Height() int64 {
	return i.BlockNumber.Int64()
}
