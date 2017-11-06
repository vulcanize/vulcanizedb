package repositories

import "github.com/8thlight/vulcanizedb/pkg/core"

type Repository interface {
	CreateBlock(block core.Block)
}
