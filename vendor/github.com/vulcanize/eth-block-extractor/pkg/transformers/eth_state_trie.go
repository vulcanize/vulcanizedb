package transformers

import (
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/vulcanize/eth-block-extractor/pkg/db"
	"github.com/vulcanize/eth-block-extractor/pkg/ipfs"
)

type EthStateTrieTransformer struct {
	database             db.Database
	stateTriePublisher   ipfs.Publisher
	storageTriePublisher ipfs.Publisher
}

func NewEthStateTrieTransformer(database db.Database, stateTriePublisher, storageTriePublisher ipfs.Publisher) *EthStateTrieTransformer {
	return &EthStateTrieTransformer{
		database:             database,
		stateTriePublisher:   stateTriePublisher,
		storageTriePublisher: storageTriePublisher,
	}
}

func (t EthStateTrieTransformer) Execute(startingBlockNumber int64, endingBlockNumber int64) error {
	if endingBlockNumber < startingBlockNumber {
		return ErrInvalidRange
	}
	for i := startingBlockNumber; i <= endingBlockNumber; i++ {
		root, err := t.getStateRootForBlock(i)
		if err != nil {
			return err
		}

		stateTrieNodes, storageTrieNodes, err := t.database.GetStateAndStorageTrieNodes(root)
		if err != nil {
			return fmt.Errorf("Error fetching state trie for block %d: %s\n", i, err)
		}

		err = t.writeStateTrieNodesToIpfs(stateTrieNodes)
		if err != nil {
			return err
		}

		err = t.writeStorageTrieNodesToIpfs(storageTrieNodes)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t EthStateTrieTransformer) getStateRootForBlock(blockNumber int64) (root common.Hash, err error) {
	header := t.database.GetBlockHeaderByBlockNumber(blockNumber)
	if err != nil {
		return root, err
	}
	return header.Root, nil
}

func (t EthStateTrieTransformer) writeStateTrieNodesToIpfs(stateTrieNodes [][]byte) error {
	for _, node := range stateTrieNodes {
		output, err := t.stateTriePublisher.Write(node)
		if err != nil {
			return fmt.Errorf("Error writing state trie node to ipfs: %s\n", err.Error())
		}
		log.Println("Created ipld: ", output)
	}
	return nil
}

func (t EthStateTrieTransformer) writeStorageTrieNodesToIpfs(storageTrieNodes [][]byte) error {
	for _, node := range storageTrieNodes {
		output, err := t.storageTriePublisher.Write(node)
		if err != nil {
			return fmt.Errorf("Error writing storage trie node to ipfs: %s\n", err.Error())
		}
		log.Println("Created ipld: ", output)
	}
	return nil
}
