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

package common

import (
	"bytes"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/makerdao/vulcanizedb/pkg/core"
)

func ToCoreReceipts(gethReceipts types.Receipts) ([]core.Receipt, error) {
	var coreReceipts []core.Receipt
	for _, receipt := range gethReceipts {
		coreReceipt, err := ToCoreReceipt(receipt)
		if err != nil {
			return nil, err
		}
		coreReceipts = append(coreReceipts, coreReceipt)
	}
	return coreReceipts, nil
}

func ToCoreReceipt(gethReceipt *types.Receipt) (core.Receipt, error) {
	bloom := hexutil.Encode(gethReceipt.Bloom.Bytes())
	var postState string
	var status int
	postState, status = postStateOrStatus(gethReceipt)
	logs := dereferenceLogs(gethReceipt)
	contractAddress := setContractAddress(gethReceipt)

	rlpBuff := new(bytes.Buffer)
	receiptForStorage := types.ReceiptForStorage(*gethReceipt)
	err := receiptForStorage.EncodeRLP(rlpBuff)
	if err != nil {
		return core.Receipt{}, err
	}
	return core.Receipt{
		Bloom:             bloom,
		ContractAddress:   contractAddress,
		CumulativeGasUsed: gethReceipt.CumulativeGasUsed,
		GasUsed:           gethReceipt.GasUsed,
		Logs:              logs,
		StateRoot:         postState,
		TxHash:            gethReceipt.TxHash.Hex(),
		Status:            status,
		Rlp:               rlpBuff.Bytes(),
	}, nil
}

func setContractAddress(gethReceipt *types.Receipt) string {
	emptyAddress := common.Address{}.Bytes()
	if bytes.Equal(gethReceipt.ContractAddress.Bytes(), emptyAddress) {
		return ""
	}
	return gethReceipt.ContractAddress.Hex()
}

func dereferenceLogs(gethReceipt *types.Receipt) []core.FullSyncLog {
	logs := []core.FullSyncLog{}
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
