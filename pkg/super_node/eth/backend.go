package eth

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/hashicorp/golang-lru"

	"github.com/vulcanize/vulcanizedb/pkg/ipfs"
	"github.com/vulcanize/vulcanizedb/pkg/super_node"
)

var (
	errPendingBlockNumber = errors.New("pending block number not supported")
)

type Backend struct {
	retriever super_node.CIDRetriever
	fetcher   ipfs.IPLDFetcher
	db        *postgres.DB

	headerCache *lru.Cache // Cache for the most recent block headers
	tdCache     *lru.Cache // Cache for the most recent block total difficulties
	numberCache *lru.Cache // Cache for the most recent block numbers
}

func NewEthBackend(r super_node.CIDRetriever, f ipfs.IPLDFetcher) *Backend {
	return &Backend{
		retriever: r,
		fetcher:   f,
		db:        r.Database(),
	}
}

func (b *Backend) HeaderByNumber(ctx context.Context, blockNumber rpc.BlockNumber) (*types.Header, error) {
	number := blockNumber.Int64()
	var err error
	if blockNumber == rpc.LatestBlockNumber {
		number, err = b.retriever.RetrieveLastBlockNumber()
		if err != nil {
			return nil, err
		}
	}
	if blockNumber == rpc.PendingBlockNumber {
		return nil, errPendingBlockNumber
	}
	// Retrieve the CIDs for headers at this height
	tx, err := b.db.Beginx()
	if err != nil {
		return nil, err
	}
	headerCids, err := b.retriever.RetrieveHeaderCIDs(tx, number)
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	// If there are none, throw an error
	if len(headerCids) < 1 {
		return nil, fmt.Errorf("header at block %d is not available", number)
	}
	// Fetch the header IPLDs for those CIDs
	headerIPLDs, err := b.fetcher.FetchHeaders(headerCids)
	if err != nil {
		return nil, err
	}
	// Decode the first header at this block height and return it
	// We throw an error in FetchHeaders() if the number of headers does not match the number of CIDs and we already
	// confirmed the number of CIDs is greater than 0 so there is no need to bound check the slice before accessing
	header := new(types.Header)
	if err := rlp.DecodeBytes(headerIPLDs[0].RawData(), header); err != nil {
		return nil, err
	}
	return header, nil
}

func (b *Backend) GetTd(blockHash common.Hash) *big.Int {
	panic("implement me")
}
