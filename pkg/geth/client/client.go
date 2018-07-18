package client

import (
	"context"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
)

type Client struct {
	client *ethclient.Client
}

func NewClient(client *ethclient.Client) Client {
	return Client{client: client}
}

func (client Client) BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error) {
	return client.client.BlockByNumber(ctx, number)
}

func (client Client) CallContract(ctx context.Context, msg ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	return client.client.CallContract(ctx, msg, blockNumber)
}

func (client Client) FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]types.Log, error) {
	return client.client.FilterLogs(ctx, q)
}

func (client Client) HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error) {
	return client.client.HeaderByNumber(ctx, number)
}
