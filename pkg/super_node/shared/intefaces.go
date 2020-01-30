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

// PayloadStreamer streams chain-specific payloads to the provided channel
type PayloadStreamer interface {
	Stream(payloadChan chan interface{}) (ClientSubscription, error)
}

// PayloadFetcher fetches chain-specific payloads
type PayloadFetcher interface {
	FetchAt(blockHeights []uint64) ([]interface{}, error)
}

// PayloadConverter converts chain-specific payloads into IPLD payloads for publishing
type PayloadConverter interface {
	Convert(payload interface{}) (interface{}, error)
}

// IPLDPublisher publishes IPLD payloads and returns a CID payload for indexing
type IPLDPublisher interface {
	Publish(payload interface{}) (interface{}, error)
}

// CIDIndexer indexes a CID payload in Postgres
type CIDIndexer interface {
	Index(cids interface{}) error
}

// ResponseFilterer applies a filter to an IPLD payload to return a subscription response packet
type ResponseFilterer interface {
	Filter(filter, payload interface{}) (response interface{}, err error)
}

// CIDRetriever retrieves cids according to a provided filter and returns a CID wrapper
type CIDRetriever interface {
	Retrieve(filter interface{}, blockNumber int64) (interface{}, bool, error)
	RetrieveFirstBlockNumber() (int64, error)
	RetrieveLastBlockNumber() (int64, error)
	RetrieveGapsInData() ([]Gap, error)
}

// IPLDFetcher uses a CID wrapper to fetch an IPLD wrapper
type IPLDFetcher interface {
	Fetch(cids interface{}) (interface{}, error)
}

// IPLDResolver resolves an IPLD wrapper into chain-specific payloads
type IPLDResolver interface {
	Resolve(iplds interface{}) (interface{}, error)
}

// ClientSubscription is a general interface for chain data subscriptions
type ClientSubscription interface {
	Err() <-chan error
	Unsubscribe()
}

// DagPutter is a general interface for a dag putter
type DagPutter interface {
	DagPut(raw interface{}) ([]string, error)
}
