package eth_state_trie

import (
	"github.com/vulcanize/eth-block-extractor/pkg/ipfs"
	"github.com/vulcanize/eth-block-extractor/pkg/ipfs/util"
)

const (
	EthStateTrieNodeCode = 0x96
)

type StateTrieDagPutter struct {
	adder ipfs.Adder
}

func NewStateTrieDagPutter(adder ipfs.Adder) *StateTrieDagPutter {
	return &StateTrieDagPutter{adder: adder}
}

func (stdp StateTrieDagPutter) DagPut(raw interface{}) ([]string, error) {
	input := raw.([]byte)
	stateTrieNode, err := stdp.getStateTrieNode(input)
	if err != nil {
		return nil, err
	}
	err = stdp.adder.Add(stateTrieNode)
	if err != nil {
		return nil, err
	}
	return []string{stateTrieNode.Cid().String()}, nil
}

func (stdp StateTrieDagPutter) getStateTrieNode(raw []byte) (*EthStateTrieNode, error) {
	stateTrieNodeCid, err := util.RawToCid(EthStateTrieNodeCode, raw)
	if err != nil {
		return nil, err
	}
	stateTrieNode := &EthStateTrieNode{
		cid:     stateTrieNodeCid,
		rawdata: raw,
	}
	return stateTrieNode, nil
}
