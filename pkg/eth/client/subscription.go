package client

import "github.com/ethereum/go-ethereum/rpc"

type Subscription struct {
	RpcSubscription *rpc.ClientSubscription
}

func (sub Subscription) Err() <-chan error {
	return sub.RpcSubscription.Err()
}

func (sub Subscription) Unsubscribe() {
	sub.RpcSubscription.Unsubscribe()
}
