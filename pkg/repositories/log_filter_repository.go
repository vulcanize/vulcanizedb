package repositories

import "github.com/vulcanize/vulcanizedb/pkg/filters"

type FilterRepository interface {
	AddFilter(filter filters.LogFilter) error
}

func (pg Postgres) AddFilter(query filters.LogFilter) error {
	_, err := pg.Db.Exec(
		`INSERT INTO log_filters 
        (name, from_block, to_block, address, topic0, topic1, topic2, topic3)
        VALUES ($1, NULLIF($2, -1), NULLIF($3, -1), $4, NULLIF($5, ''), NULLIF($6, ''), NULLIF($7, ''), NULLIF($8, ''))`,
		query.Name, query.FromBlock, query.ToBlock, query.Address, query.Topics[0], query.Topics[1], query.Topics[2], query.Topics[3])
	if err != nil {
		return err
	}
	return nil
}
