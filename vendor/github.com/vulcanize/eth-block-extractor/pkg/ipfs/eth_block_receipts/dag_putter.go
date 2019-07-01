package eth_block_receipts

import (
	"bytes"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ipfs/go-cid"

	"github.com/vulcanize/eth-block-extractor/pkg/ipfs"
	"github.com/vulcanize/eth-block-extractor/pkg/ipfs/util"
)

type EthBlockReceiptDagPutter struct {
	adder ipfs.Adder
}

func NewEthBlockReceiptDagPutter(adder ipfs.Adder) *EthBlockReceiptDagPutter {
	return &EthBlockReceiptDagPutter{adder: adder}
}

func (dagPutter *EthBlockReceiptDagPutter) DagPut(raw interface{}) ([]string, error) {
	input := raw.(types.Receipts)
	var output []string
	for _, r := range input {
		receiptForStorage := (*types.ReceiptForStorage)(r)
		node, err := getReceiptNode(receiptForStorage)
		if err != nil {
			return nil, err
		}
		err = dagPutter.adder.Add(node)
		if err != nil {
			return nil, err
		}
		output = append(output, node.cid.String())
	}
	return output, nil
}

func getReceiptNode(receipt *types.ReceiptForStorage) (*EthReceiptNode, error) {
	buffer := new(bytes.Buffer)
	err := receipt.EncodeRLP(buffer)
	if err != nil {
		return nil, err
	}
	receiptCid, err := util.RawToCid(cid.EthTxReceipt, buffer.Bytes())
	if err != nil {
		return nil, err
	}
	node := &EthReceiptNode{
		raw: buffer.Bytes(),
		cid: receiptCid,
	}
	return node, nil
}
