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

package shared

import (
	"math/big"
)

// PayloadStreamer streams chain-specific payloads to the provided channel
type PayloadStreamer interface {
	Stream(payloadChan chan RawChainData) (ClientSubscription, error)
}

// PayloadFetcher fetches chain-specific payloads
type PayloadFetcher interface {
	FetchAt(blockHeights []uint64) ([]RawChainData, error)
}

// PayloadConverter converts chain-specific payloads into IPLD payloads for publishing
type PayloadConverter interface {
	Convert(payload RawChainData) (ConvertedData, error)
}

// IPLDPublisher publishes IPLD payloads and returns a CID payload for indexing
type IPLDPublisher interface {
	Publish(payload ConvertedData) (CIDsForIndexing, error)
}

// CIDIndexer indexes a CID payload in Postgres
type CIDIndexer interface {
	Index(cids CIDsForIndexing) error
}

// ResponseFilterer applies a filter to an IPLD payload to return a subscription response packet
type ResponseFilterer interface {
	Filter(filter SubscriptionSettings, payload ConvertedData) (response IPLDs, err error)
}

// CIDRetriever retrieves cids according to a provided filter and returns a CID wrapper
type CIDRetriever interface {
	Retrieve(filter SubscriptionSettings, blockNumber int64) ([]CIDsForFetching, bool, error)
	RetrieveFirstBlockNumber() (int64, error)
	RetrieveLastBlockNumber() (int64, error)
	RetrieveGapsInData(validationLevel int) ([]Gap, error)
}

// IPLDFetcher uses a CID wrapper to fetch an IPLD wrapper
type IPLDFetcher interface {
	Fetch(cids CIDsForFetching) (IPLDs, error)
}

// ClientSubscription is a general interface for chain data subscriptions
type ClientSubscription interface {
	Err() <-chan error
	Unsubscribe()
}

// Cleaner is for cleaning out data from the cache within the given ranges
type Cleaner interface {
	Clean(rngs [][2]uint64, t DataType) error
	ResetValidation(rngs [][2]uint64) error
}

// Validator is for validating sections of data using chain-specific procedures
type Validator interface {
	Validate(errChan chan error) error
}

// SubscriptionSettings is the interface every subscription filter type needs to satisfy, no matter the chain
// Further specifics of the underlying filter type depend on the internal needs of the types
// which satisfy the ResponseFilterer and CIDRetriever interfaces for a specific chain
// The underlying type needs to be rlp serializable
type SubscriptionSettings interface {
	StartingBlock() *big.Int
	EndingBlock() *big.Int
	ChainType() ChainType
	HistoricalData() bool
	HistoricalDataOnly() bool
}
