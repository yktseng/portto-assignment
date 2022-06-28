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
}

func (r *RPC) Connect(endpoint string) bool {
	client, err := ethclient.Dial(endpoint)
	if err != nil {
		return false
	}
	r.client = client
	return true
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
