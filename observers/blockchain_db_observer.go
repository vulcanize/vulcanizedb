package observers

import (
	"github.com/8thlight/vulcanizedb/core"
	"github.com/8thlight/vulcanizedb/repositories"
	"github.com/jmoiron/sqlx"
)

type BlockchainDBObserver struct {
	Db *sqlx.DB
}

func (observer BlockchainDBObserver) NotifyBlockAdded(block core.Block) {
	repositories.NewPostgres(observer.Db).CreateBlock(block)
}
