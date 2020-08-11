package mocks

import (
	"github.com/ethereum/go-ethereum/eth/filters"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
)

type MockStoragediffStreamer struct {
	subscribeError     error
	ClientSubscription *fakes.MockSubscription
	PassedPayloadChan  chan filters.Payload
	streamPayloads     []filters.Payload
}

func (streamer *MockStoragediffStreamer) Stream(statediffPayloadChan chan filters.Payload) (core.Subscription, error) {
	streamer.PassedPayloadChan = statediffPayloadChan

	go func() {
		for _, payload := range streamer.streamPayloads {
			streamer.PassedPayloadChan <- payload
		}
	}()

	return streamer.ClientSubscription, streamer.subscribeError
}

func (streamer *MockStoragediffStreamer) SetSubscribeError(err error) {
	streamer.subscribeError = err
}

func (streamer *MockStoragediffStreamer) SetPayloads(payloads []filters.Payload) {
	streamer.streamPayloads = payloads
}
