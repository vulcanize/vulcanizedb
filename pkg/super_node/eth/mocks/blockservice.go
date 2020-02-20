package mocks

import (
	"context"
	"errors"

	"github.com/ipfs/go-block-format"
	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-ipfs-blockstore"
	"github.com/ipfs/go-ipfs-exchange-interface"
)

// MockIPFSBlockService is a mock for testing the ipfs fetcher
type MockIPFSBlockService struct {
	Blocks map[cid.Cid]blocks.Block
}

// GetBlock is used to retrieve a block from the mock BlockService
func (bs *MockIPFSBlockService) GetBlock(ctx context.Context, c cid.Cid) (blocks.Block, error) {
	if bs.Blocks == nil {
		return nil, errors.New("BlockService has not been initialized")
	}
	blk, ok := bs.Blocks[c]
	if ok {
		return blk, nil
	}
	return nil, nil
}

// GetBlocks is used to retrieve a set of blocks from the mock BlockService
func (bs *MockIPFSBlockService) GetBlocks(ctx context.Context, cs []cid.Cid) <-chan blocks.Block {
	if bs.Blocks == nil {
		panic("BlockService has not been initialized")
	}
	blkChan := make(chan blocks.Block)
	go func() {
		for _, c := range cs {
			blk, ok := bs.Blocks[c]
			if ok {
				blkChan <- blk
			}
		}
		close(blkChan)
	}()
	return blkChan
}

// AddBlock adds a block to the mock BlockService
func (bs *MockIPFSBlockService) AddBlock(blk blocks.Block) error {
	if bs.Blocks == nil {
		bs.Blocks = make(map[cid.Cid]blocks.Block)
	}
	bs.Blocks[blk.Cid()] = blk
	return nil
}

// AddBlocks adds a set of blocks to the mock BlockService
func (bs *MockIPFSBlockService) AddBlocks(blks []blocks.Block) error {
	if bs.Blocks == nil {
		bs.Blocks = make(map[cid.Cid]blocks.Block)
	}
	for _, block := range blks {
		bs.Blocks[block.Cid()] = block
	}
	return nil
}

// Close is here to satisfy the interface
func (*MockIPFSBlockService) Close() error {
	panic("implement me")
}

// Blockstore is here to satisfy the interface
func (*MockIPFSBlockService) Blockstore() blockstore.Blockstore {
	panic("implement me")
}

// DeleteBlock is here to satisfy the interface
func (*MockIPFSBlockService) DeleteBlock(c cid.Cid) error {
	panic("implement me")
}

// Exchange is here to satisfy the interface
func (*MockIPFSBlockService) Exchange() exchange.Interface {
	panic("implement me")
}
