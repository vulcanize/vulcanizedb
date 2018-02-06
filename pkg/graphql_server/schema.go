package graphql_server

import (
	"github.com/vulcanize/vulcanizedb/pkg/core"
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
        watchedEvents(name: String!): WatchedEventList
	}

    type LogFilter {
        name: String!
        fromBlock: Int
        toBlock: Int
        address: String!
        topics: [String]!
    }

  	type WatchedEventList{
		total: Int!
		watchedEvents: [WatchedEvent]!
	} 

    type WatchedEvent {
        name: String!
        blockNumber: Int!
        address: String!
        tx_hash: String!
        topic0: String!
        topic1: String!
        topic2: String!
        topic3: String!
        data: String!
    }
`

type Resolver struct {
	repository repositories.Repository
}

func NewResolver(repository repositories.Repository) *Resolver {
	return &Resolver{repository: repository}
}

func (r *Resolver) WatchedEvents(args struct {
	Name string
}) (*watchedEventsResolver, error) {
	wel, err := r.repository.GetWatchedEvents(args.Name)
	if err != nil {
		return &watchedEventsResolver{}, err
	}
	return &watchedEventsResolver{wels: wel}, err
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

type logFilter struct {
	*filters.LogFilter
}

type logFilterResolver struct {
	lf *logFilter
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

type watchedEventsResolver struct {
	wels []*core.WatchedEvent
}

func (wels watchedEventsResolver) WatchedEvents() []*watchedEventResolver {
	return newWatchedEventsResolver(wels.wels)
}

func (wels watchedEventsResolver) Total() int32 {
	return int32(len(wels.wels))
}

func newWatchedEventsResolver(wel []*core.WatchedEvent) []*watchedEventResolver {
	watchedEventResolvers := make([]*watchedEventResolver, 0)
	for _, e := range wel {
		watchedEventResolvers = append(watchedEventResolvers, &watchedEventResolver{&watchedEvent{e}})
	}
	return watchedEventResolvers
}

type watchedEvent struct {
	*core.WatchedEvent
}

type watchedEventResolver struct {
	we *watchedEvent
}

func (wr watchedEventResolver) Name() string {
	return wr.we.Name
}

func (wr watchedEventResolver) BlockNumber() int32 {
	return int32(wr.we.BlockNumber)
}

func (wr watchedEventResolver) Address() string {
	return wr.we.Address
}

func (wr watchedEventResolver) TxHash() string {
	return wr.we.TxHash
}

func (wr watchedEventResolver) Topic0() string {
	return wr.we.Topic0
}

func (wr watchedEventResolver) Topic1() string {
	return wr.we.Topic1
}

func (wr watchedEventResolver) Topic2() string {
	return wr.we.Topic2
}

func (wr watchedEventResolver) Topic3() string {
	return wr.we.Topic3
}

func (wr watchedEventResolver) Data() string {
	return wr.we.Data
}
