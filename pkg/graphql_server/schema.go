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

func (r *Resolver) LogFilter(args struct {
	Name string
}) (*logFilterResolver, error) {
	logFilter, err := r.repository.GetFilter(args.Name)
	if err != nil {
		return &logFilterResolver{}, err
	}
	return &logFilterResolver{&logFilter}, nil
}

type logFilterResolver struct {
	lf *filters.LogFilter
}

func (lfr *logFilterResolver) Name() string {
	return lfr.lf.Name
}

func (lfr *logFilterResolver) FromBlock() *int32 {
	fromBlock := int32(lfr.lf.FromBlock)
	return &fromBlock
}

func (lfr *logFilterResolver) ToBlock() *int32 {
	toBlock := int32(lfr.lf.ToBlock)
	return &toBlock
}

func (lfr *logFilterResolver) Address() string {
	return lfr.lf.Address
}

func (lfr *logFilterResolver) Topics() []*string {
	var topics = make([]*string, 4)
	for i := range topics {
		if lfr.lf.Topics[i] != "" {
			topics[i] = &lfr.lf.Topics[i]
		}
	}
	return topics
}

func (r *Resolver) WatchedEvents(args struct {
	Name string
}) (*watchedEventsResolver, error) {
	watchedEvents, err := r.repository.GetWatchedEvents(args.Name)
	if err != nil {
		return &watchedEventsResolver{}, err
	}
	return &watchedEventsResolver{watchedEvents: watchedEvents}, err
}

type watchedEventsResolver struct {
	watchedEvents []*core.WatchedEvent
}

func (wesr watchedEventsResolver) WatchedEvents() []*watchedEventResolver {
	return resolveWatchedEvents(wesr.watchedEvents)
}

func (wesr watchedEventsResolver) Total() int32 {
	return int32(len(wesr.watchedEvents))
}

func resolveWatchedEvents(watchedEvents []*core.WatchedEvent) []*watchedEventResolver {
	watchedEventResolvers := make([]*watchedEventResolver, 0)
	for _, watchedEvent := range watchedEvents {
		watchedEventResolvers = append(watchedEventResolvers, &watchedEventResolver{watchedEvent})
	}
	return watchedEventResolvers
}

type watchedEventResolver struct {
	we *core.WatchedEvent
}

func (wer watchedEventResolver) Name() string {
	return wer.we.Name
}

func (wer watchedEventResolver) BlockNumber() int32 {
	return int32(wer.we.BlockNumber)
}

func (wer watchedEventResolver) Address() string {
	return wer.we.Address
}

func (wer watchedEventResolver) TxHash() string {
	return wer.we.TxHash
}

func (wer watchedEventResolver) Topic0() string {
	return wer.we.Topic0
}

func (wer watchedEventResolver) Topic1() string {
	return wer.we.Topic1
}

func (wer watchedEventResolver) Topic2() string {
	return wer.we.Topic2
}

func (wer watchedEventResolver) Topic3() string {
	return wer.we.Topic3
}

func (wer watchedEventResolver) Data() string {
	return wer.we.Data
}
