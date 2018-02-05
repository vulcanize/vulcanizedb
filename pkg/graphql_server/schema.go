package graphql_server

import (
	"github.com/vulcanize/vulcanizedb/pkg/filters"
	"github.com/vulcanize/vulcanizedb/pkg/repositories"
)

var Schema = `
	schema {
		query: Query
	}
	# The query type, represents all of the entry points into our object graph
	type Query {
		block(number: Int!): Block
		logs(address: String!, blockNumber: Int!): [Log]
	}
    type LogFilter {
        name: String!
        fromBlock: Int
        toBlock: Int
        address: String!
        topic0: String
        topic1: String
        topic2: String
        topic3: String
    }

    type WatchedEvent {
        name: String!
        blockNumber: Int 
        address: String!
        txHash: String
        topic0: String
        topic1: String
        topic2: String
        topic3: String        
        data: String
    }
`

type Resolver struct {
	repository repositories.Repository
}

func NewResolver(repository repositories.Repository) *Resolver {
	return &Resolver{repository: repository}
}

type LogFilter struct {
	*filters.LogFilter
}

type logFilterResolver struct {
	lf *LogFilter
}

func (r *Resolver) LogFilter(args struct {
	Name string
}) *logFilterResolver {
	//logFilter := r.repository.GetLogFilter(args.Name)
	return &logFilterResolver{logFilter}
}

func (lr *logFilterResolver) Name() string {
	return lr.lf.Name
}

func (lr *logFilterResolver) FromBlock() int32 {
	return int32(lr.lf.FromBlock)
}

func (lr *logFilterResolver) ToBlock() int32 {
	return int32(lr.lf.ToBlock)
}
