package repositories

import "github.com/8thlight/vulcanizedb/core"

type Repository interface {
	CreateBlock(block core.Block)
}
