package repositories

import (
	"errors"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/vulcanize/vulcanizedb/pkg/config"
	"github.com/vulcanize/vulcanizedb/pkg/core"
)

type Postgres struct {
	Db     *sqlx.DB
	node   core.Node
	nodeId int64
}

var (
	ErrDBInsertFailed     = errors.New("postgres: insert failed")
	ErrDBDeleteFailed     = errors.New("postgres: delete failed")
	ErrDBConnectionFailed = errors.New("postgres: db connection failed")
	ErrUnableToSetNode    = errors.New("postgres: unable to set node")
)

func NewPostgres(databaseConfig config.Database, node core.Node) (*Postgres, error) {
	connectString := config.DbConnectionString(databaseConfig)
	db, err := sqlx.Connect("postgres", connectString)
	if err != nil {
		return &Postgres{}, ErrDBConnectionFailed
	}
	pg := Postgres{Db: db, node: node}
	err = pg.CreateNode(&node)
	if err != nil {
		return &Postgres{}, ErrUnableToSetNode
	}
	return &pg, nil
}
