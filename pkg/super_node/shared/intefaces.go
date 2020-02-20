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
	"github.com/ethereum/go-ethereum/rpc"
)

// ResponseFilterer applies a filter to the streamed payload and returns a subscription response packet
type ResponseFilterer interface {
	Filter(filter, payload interface{}) (response interface{}, err error)
}

// CIDIndexer indexes a set of cids with their associated meta data in Postgres
type CIDIndexer interface {
	Index(cids interface{}) error
}

// CIDRetriever retrieves cids according to a provided filter and returns a cid
type CIDRetriever interface {
	Retrieve(filter interface{}, blockNumber int64) (interface{}, bool, error)
	RetrieveFirstBlockNumber() (int64, error)
	RetrieveLastBlockNumber() (int64, error)
	RetrieveGapsInData() ([]Gap, error)
}

type PayloadStreamer interface {
	Stream(payloadChan chan interface{}) (*rpc.ClientSubscription, error)
}

type PayloadFetcher interface {
	FetchAt(blockHeights []uint64) ([]interface{}, error)
}

type IPLDFetcher interface {
	Fetch(cids interface{}) (interface{}, error)
}

type PayloadConverter interface {
	Convert(payload interface{}) (interface{}, error)
}

type IPLDPublisher interface {
	Publish(payload interface{}) (interface{}, error)
}

type IPLDResolver interface {
	Resolve(iplds interface{}) (interface{}, error)
}
