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
        logFilter(name: String!): LogFilter
	}

    type LogFilter {
        name: String!
        fromBlock: Int
        toBlock: Int
        address: String!
        topics: [String]!
    }
`

type Resolver struct {
	repository repositories.Repository
}

func NewResolver(repository repositories.Repository) *Resolver {
	return &Resolver{repository: repository}
}

type logFilter struct {
	*filters.LogFilter
}

type logFilterResolver struct {
	lf *logFilter
}

func (r *Resolver) LogFilter(args struct {
	Name string
}) (*logFilterResolver, error) {
	lf, err := r.repository.GetFilter(args.Name)
	if err != nil {
		return &logFilterResolver{}, err
	}
	return &logFilterResolver{&logFilter{&lf}}, nil
}

func (lr *logFilterResolver) Name() string {
	return lr.lf.Name
}

func (lr *logFilterResolver) FromBlock() *int32 {
	fromBlock := int32(lr.lf.FromBlock)
	return &fromBlock
}

func (lr *logFilterResolver) ToBlock() *int32 {
	toBlock := int32(lr.lf.ToBlock)
	return &toBlock
}

func (lr *logFilterResolver) Address() string {
	return lr.lf.Address
}

func (lr *logFilterResolver) Topics() []*string {
	var topics = make([]*string, 4)
	for i := range topics {
		if lr.lf.Topics[i] != "" {
			topics[i] = &lr.lf.Topics[i]
		}
	}
	return topics
}
