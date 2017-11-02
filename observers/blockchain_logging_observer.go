package observers

import (
	"fmt"
	"time"

	"github.com/8thlight/vulcanizedb/core"
)

type BlockchainLoggingObserver struct{}

func (blockchainObserver BlockchainLoggingObserver) NotifyBlockAdded(block core.Block) {
	fmt.Printf("New block was added: %d\n"+
		"\tTime: %v\n"+
		"\tGas Limit: %d\n"+
		"\tGas Used: %d\n"+
		"\tNumber of Transactions %d\n", block.Number, time.Unix(block.Time, 0), block.GasLimit, block.GasUsed, len(block.Transactions))
}
