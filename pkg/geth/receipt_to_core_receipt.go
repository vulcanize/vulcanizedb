package geth

import (
	"math/big"

	"bytes"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/vulcanizedb/pkg/core"
)

func BigTo64(n *big.Int) int64 {
	if n != nil {
		return n.Int64()
	}
	return 0
}

func ReceiptToCoreReceipt(gethReceipt *types.Receipt) core.Receipt {
	bloom := hexutil.Encode(gethReceipt.Bloom.Bytes())
	var postState string
	var status int
	postState, status = postStateOrStatus(gethReceipt)
	logs := dereferenceLogs(gethReceipt)
	contractAddress := setContractAddress(gethReceipt)

	return core.Receipt{
		Bloom:             bloom,
		ContractAddress:   contractAddress,
		CumulativeGasUsed: gethReceipt.CumulativeGasUsed.Int64(),
		GasUsed:           gethReceipt.GasUsed.Int64(),
		Logs:              logs,
		StateRoot:         postState,
		TxHash:            gethReceipt.TxHash.Hex(),
		Status:            status,
	}
}

func setContractAddress(gethReceipt *types.Receipt) string {
	emptyAddress := common.Address{}.Bytes()
	if bytes.Equal(gethReceipt.ContractAddress.Bytes(), emptyAddress) {
		return ""
	}
	return gethReceipt.ContractAddress.Hex()
}

func dereferenceLogs(gethReceipt *types.Receipt) []core.Log {
	logs := []core.Log{}
	for _, log := range gethReceipt.Logs {
		logs = append(logs, ToCoreLog(*log))
	}
	return logs
}

func postStateOrStatus(gethReceipts *types.Receipt) (string, int) {
	if len(gethReceipts.PostState) != 0 {
		return hexutil.Encode(gethReceipts.PostState), -99
	}
	return "", int(gethReceipts.Status)
}
