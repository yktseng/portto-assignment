package myeth

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/yktseng/portto-assignment/internal/block"
)

type RPC struct {
	client *ethclient.Client
	endpoint string
}

func (r *RPC) Connect(endpoint string) bool {
	client, err := ethclient.Dial(endpoint)
	if err != nil {
		return false
	}
	r.client = client
	r.endpoint = endpoint
	return true
}

func (r *RPC) Endpoint() string {
	return r.endpoint
}

func (r *RPC) GetBlock(ctx context.Context,
	num *big.Int) (*block.Block, error) {
	raw, err := r.client.BlockByNumber(ctx, num)
	if err != nil {
		return nil, err
	}
	block, err := block.FromGethBlock(raw)
	if err != nil {
		return nil, err
	}
	block.Stable = true
	return block, nil
}

func (r *RPC) GetTx(ctx context.Context,
	hash common.Hash) (*types.Transaction, error) {
	tx, _, err := r.client.TransactionByHash(ctx, hash)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (r *RPC) GetTxReceipt(ctx context.Context,
	hash common.Hash) (*types.Receipt, error) {
	receipt, err := r.client.TransactionReceipt(ctx, hash)
	if err != nil {
		return nil, err
	}
	return receipt, nil
}

func(r *RPC) GetNewestBlock(ctx context.Context) (*big.Int, error) {
	b, err := r.client.BlockByNumber(ctx, nil)
	if err != nil {
		return nil, err
	}
	return b.Number(), nil
}
