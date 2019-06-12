package sqlds

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	ds "github.com/ipfs/go-datastore"
	dsq "github.com/ipfs/go-datastore/query"
)

type Queries interface {
	Delete() string
	Exists() string
	Get() string
	Put() string
	Query() string
	Prefix() string
	Limit() string
	Offset() string
	GetSize() string
}

type Datastore struct {
	db      *sql.DB
	queries Queries
}

// NewDatastore returns a new datastore
func NewDatastore(db *sql.DB, queries Queries) *Datastore {
	return &Datastore{db: db, queries: queries}
}

type batch struct {
	db      *sql.DB
	queries Queries
	txn     *sql.Tx
}

func (b *batch) GetTransaction() (*sql.Tx, error) {
	if b.txn != nil {
		return b.txn, nil
	}

	newTransaction, err := b.db.Begin()
	if err != nil {
		if newTransaction != nil {
			newTransaction.Rollback()
		}

		return nil, err
	}

	b.txn = newTransaction
	return newTransaction, nil
}

func (b *batch) Put(key ds.Key, val []byte) error {
	if val == nil {
		return ds.ErrInvalidType
	}

	txn, err := b.GetTransaction()
	if err != nil {
		b.txn.Rollback()
		return err
	}

	_, err = txn.Exec(b.queries.Put(), key.String(), val)
	if err != nil {
		b.txn.Rollback()
		return err
	}

	return nil
}

func (b *batch) Delete(key ds.Key) error {
	txn, err := b.GetTransaction()
	if err != nil {
		b.txn.Rollback()
	}

	_, err = txn.Exec(b.queries.Delete(), key.String())
	if err != nil {
		b.txn.Rollback()
		return err
	}

	return err
}

func (b *batch) Commit() error {
	if b.txn == nil {
		return errors.New("no transaction started, cannot commit")
	}
	var err = b.txn.Commit()
	if err != nil {
		b.txn.Rollback()
		return err
	}

	return nil
}

func (d *Datastore) Batch() (ds.Batch, error) {
	batch := &batch{
		db:      d.db,
		queries: d.queries,
		txn:     nil,
	}

	return batch, nil
}

func (d *Datastore) Close() error {
	return d.db.Close()
}

func (d *Datastore) Delete(key ds.Key) error {
	result, err := d.db.Exec(d.queries.Delete(), key.String())
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ds.ErrNotFound
	}

	return nil
}

func (d *Datastore) Get(key ds.Key) (value []byte, err error) {
	row := d.db.QueryRow(d.queries.Get(), key.String())
	var out []byte

	switch err := row.Scan(&out); err {
	case sql.ErrNoRows:
		return nil, ds.ErrNotFound
	case nil:
		return out, nil
	default:
		return nil, err
	}
}

func (d *Datastore) Has(key ds.Key) (exists bool, err error) {
	row := d.db.QueryRow(d.queries.Exists(), key.String())

	switch err := row.Scan(&exists); err {
	case sql.ErrNoRows:
		return exists, nil
	case nil:
		return exists, nil
	default:
		return exists, err
	}
}

func (d *Datastore) Put(key ds.Key, value []byte) error {
	if value == nil {
		return ds.ErrInvalidType
	}

	_, err := d.db.Exec(d.queries.Put(), key.String(), value)
	if err != nil {
		return err
	}

	return nil
}

func (d *Datastore) Query(q dsq.Query) (dsq.Results, error) {
	raw, err := d.RawQuery(q)
	if err != nil {
		return nil, err
	}

	for _, f := range q.Filters {
		raw = dsq.NaiveFilter(raw, f)
	}

	raw = dsq.NaiveOrder(raw, q.Orders...)

	return raw, nil
}

func (d *Datastore) RawQuery(q dsq.Query) (dsq.Results, error) {
	var rows *sql.Rows
	var err error

	if q.Prefix != "" {
		rows, err = QueryWithParams(d, q)
	} else {
		rows, err = d.db.Query(d.queries.Query())
	}

	if err != nil {
		return nil, err
	}

	var entries []dsq.Entry
	defer rows.Close()

	for rows.Next() {
		var key string
		var out []byte
		err := rows.Scan(&key, &out)

		if err != nil {
			log.Fatal("Error reading rows from query")
		}

		entry := dsq.Entry{
			Key:   key,
			Value: out,
		}

		entries = append(entries, entry)
	}

	results := dsq.ResultsWithEntries(q, entries)
	return results, nil
}

func (d *Datastore) GetSize(key ds.Key) (int, error) {
	row := d.db.QueryRow(d.queries.GetSize(), key.String())
	var size int

	switch err := row.Scan(&size); err {
	case sql.ErrNoRows:
		return -1, ds.ErrNotFound
	case nil:
		return size, nil
	default:
		return 0, err
	}
}

// QueryWithParams applies prefix, limit, and offset params in pg query
func QueryWithParams(d *Datastore, q dsq.Query) (*sql.Rows, error) {
	var qNew = d.queries.Query()

	if q.Prefix != "" {
		qNew += fmt.Sprintf(d.queries.Prefix(), q.Prefix)
	}

	if q.Limit != 0 {
		qNew += fmt.Sprintf(d.queries.Limit(), q.Limit)
	}

	if q.Offset != 0 {
		qNew += fmt.Sprintf(d.queries.Offset(), q.Offset)
	}

	return d.db.Query(qNew)

}

var _ ds.Datastore = (*Datastore)(nil)
