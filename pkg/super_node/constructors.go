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

package super_node

import (
	"fmt"
	"github.com/btcsuite/btcd/chaincfg"

	"github.com/btcsuite/btcd/rpcclient"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/vulcanize/vulcanizedb/pkg/eth/core"
	"github.com/vulcanize/vulcanizedb/pkg/eth/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/super_node/btc"
	"github.com/vulcanize/vulcanizedb/pkg/super_node/eth"
	"github.com/vulcanize/vulcanizedb/pkg/super_node/shared"
)

// NewResponseFilterer constructs a ResponseFilterer for the provided chain type
func NewResponseFilterer(chain shared.ChainType) (shared.ResponseFilterer, error) {
	switch chain {
	case shared.Ethereum:
		return eth.NewResponseFilterer(), nil
	default:
		return nil, fmt.Errorf("invalid chain %T for filterer constructor", chain)
	}
}

// NewCIDIndexer constructs a CIDIndexer for the provided chain type
func NewCIDIndexer(chain shared.ChainType, db *postgres.DB) (shared.CIDIndexer, error) {
	switch chain {
	case shared.Ethereum:
		return eth.NewCIDIndexer(db), nil
	case shared.Bitcoin:
		return btc.NewCIDIndexer(db), nil
	default:
		return nil, fmt.Errorf("invalid chain %T for indexer constructor", chain)
	}
}

// NewCIDRetriever constructs a CIDRetriever for the provided chain type
func NewCIDRetriever(chain shared.ChainType, db *postgres.DB) (shared.CIDRetriever, error) {
	switch chain {
	case shared.Ethereum:
		return eth.NewCIDRetriever(db), nil
	default:
		return nil, fmt.Errorf("invalid chain %T for retriever constructor", chain)
	}
}

// NewPayloadStreamer constructs a PayloadStreamer for the provided chain type
func NewPayloadStreamer(chain shared.ChainType, client interface{}) (shared.PayloadStreamer, chan shared.RawChainData, error) {
	switch chain {
	case shared.Ethereum:
		ethClient, ok := client.(core.RPCClient)
		if !ok {
			var expectedClientType core.RPCClient
			return nil, nil, fmt.Errorf("ethereum payload streamer constructor expected client type %T got %T", expectedClientType, client)
		}
		streamChan := make(chan shared.RawChainData, eth.PayloadChanBufferSize)
		return eth.NewPayloadStreamer(ethClient), streamChan, nil
	case shared.Bitcoin:
		btcClientConn, ok := client.(*rpcclient.ConnConfig)
		if !ok {
			return nil, nil, fmt.Errorf("bitcoin payload streamer constructor expected client shared type %T got %T", rpcclient.ConnConfig{}, client)
		}
		streamChan := make(chan shared.RawChainData, btc.PayloadChanBufferSize)
		return btc.NewPayloadStreamer(btcClientConn), streamChan, nil
	default:
		return nil, nil, fmt.Errorf("invalid chain %T for streamer constructor", chain)
	}
}

// NewPaylaodFetcher constructs a PayloadFetcher for the provided chain type
func NewPaylaodFetcher(chain shared.ChainType, client interface{}) (shared.PayloadFetcher, error) {
	switch chain {
	case shared.Ethereum:
		batchClient, ok := client.(eth.BatchClient)
		if !ok {
			var expectedClient eth.BatchClient
			return nil, fmt.Errorf("ethereum fetcher constructor expected client type %T got %T", expectedClient, client)
		}
		return eth.NewPayloadFetcher(batchClient), nil
	default:
		return nil, fmt.Errorf("invalid chain %T for fetcher constructor", chain)
	}
}

// NewPayloadConverter constructs a PayloadConverter for the provided chain type
func NewPayloadConverter(chain shared.ChainType) (shared.PayloadConverter, error) {
	switch chain {
	case shared.Ethereum:
		return eth.NewPayloadConverter(params.MainnetChainConfig), nil
	case shared.Bitcoin:
		return btc.NewPayloadConverter(&chaincfg.MainNetParams), nil
	default:
		return nil, fmt.Errorf("invalid chain %T for converter constructor", chain)
	}
}

// NewIPLDFetcher constructs an IPLDFetcher for the provided chain type
func NewIPLDFetcher(chain shared.ChainType, ipfsPath string) (shared.IPLDFetcher, error) {
	switch chain {
	case shared.Ethereum:
		return eth.NewIPLDFetcher(ipfsPath)
	default:
		return nil, fmt.Errorf("invalid chain %T for fetcher constructor", chain)
	}
}

// NewIPLDPublisher constructs an IPLDPublisher for the provided chain type
func NewIPLDPublisher(chain shared.ChainType, ipfsPath string) (shared.IPLDPublisher, error) {
	switch chain {
	case shared.Ethereum:
		return eth.NewIPLDPublisher(ipfsPath)
	case shared.Bitcoin:
		return btc.NewIPLDPublisher(ipfsPath)
	default:
		return nil, fmt.Errorf("invalid chain %T for publisher constructor", chain)
	}
}

// NewIPLDResolver constructs an IPLDResolver for the provided chain type
func NewIPLDResolver(chain shared.ChainType) (shared.IPLDResolver, error) {
	switch chain {
	case shared.Ethereum:
		return eth.NewIPLDResolver(), nil
	default:
		return nil, fmt.Errorf("invalid chain %T for resolver constructor", chain)
	}
}

// NewPublicAPI constructs a PublicAPI for the provided chain type
func NewPublicAPI(chain shared.ChainType, db *postgres.DB, ipfsPath string) (rpc.API, error) {
	switch chain {
	case shared.Ethereum:
		backend, err := eth.NewEthBackend(db, ipfsPath)
		if err != nil {
			return rpc.API{}, err
		}
		return rpc.API{
			Namespace: eth.APIName,
			Version:   eth.APIVersion,
			Service:   eth.NewPublicEthAPI(backend),
			Public:    true,
		}, nil
	default:
		return rpc.API{}, fmt.Errorf("invalid chain %T for public api constructor", chain)
	}
}
