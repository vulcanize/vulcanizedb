package observers

import (
	"os"
	"text/template"

	"time"

	"github.com/8thlight/vulcanizedb/pkg/core"
)

const blockAddedTemplate = `
    New block was added: {{.Number}}
    Time: {{.Time | unix_time}}
    Gas Limit: {{.GasLimit}}
    Gas Used: {{.GasUsed}}
    Number of Transactions {{.Transactions | len}}

`

var funcMap = template.FuncMap{
	"unix_time": func(n int64) time.Time {
		return time.Unix(n, 0)
	},
}
var tmp = template.Must(template.New("window").Funcs(funcMap).Parse(blockAddedTemplate))

type BlockchainLoggingObserver struct{}

func (blockchainObserver BlockchainLoggingObserver) NotifyBlockAdded(block core.Block) {
	tmp.Execute(os.Stdout, block)
}
