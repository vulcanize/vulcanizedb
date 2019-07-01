package postgres

import (
	"database/sql"
	"fmt"

	"github.com/whyrusleeping/sql-datastore"

	_ "github.com/lib/pq" //postgres driver
)

// Options are the postgres datastore options, reexported here for convenience.
type Options struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

type Queries struct {
}

func (Queries) Delete() string {
	return `DELETE FROM blocks WHERE key = $1`
}

func (Queries) Exists() string {
	return `SELECT exists(SELECT 1 FROM blocks WHERE key=$1)`
}

func (Queries) Get() string {
	return `SELECT data FROM blocks WHERE key = $1`
}

func (Queries) Put() string {
	return `INSERT INTO blocks (key, data) SELECT $1, $2 WHERE NOT EXISTS ( SELECT key FROM blocks WHERE key = $1)`
}

func (Queries) Query() string {
	return `SELECT key, data FROM blocks`
}

func (Queries) Prefix() string {
	return ` WHERE key LIKE '%s%%' ORDER BY key`
}

func (Queries) Limit() string {
	return ` LIMIT %d`
}

func (Queries) Offset() string {
	return ` OFFSET %d`
}

func (Queries) GetSize() string {
	return `SELECT octet_length(data) FROM blocks WHERE key = $1`
}

// Create returns a datastore connected to postgres
func (opts *Options) Create() (*sqlds.Datastore, error) {
	opts.setDefaults()
	fmtstr := "postgresql:///%s?host=%s&port=%s&user=%s&password=%s&sslmode=disable"
	constr := fmt.Sprintf(fmtstr, opts.Database, opts.Host, opts.Port, opts.User, opts.Password)
	db, err := sql.Open("postgres", constr)
	if err != nil {
		return nil, err
	}

	return sqlds.NewDatastore(db, Queries{}), nil
}

func (opts *Options) setDefaults() {
	if opts.Host == "" {
		opts.Host = "127.0.0.1"
	}

	if opts.Port == "" {
		opts.Port = "5432"
	}

	if opts.User == "" {
		opts.User = "postgres"
	}

	if opts.Database == "" {
		opts.Database = "datastore"
	}
}
