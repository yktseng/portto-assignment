package myeth

import (
	"context"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type BlockNumHash struct {
	Hash common.Hash
	Num  *big.Int
}

func (r *RPC) SubscribeNewHeaders(ctx context.Context, output chan BlockNumHash) error {
	defer func() {
		log.Println("ws endpoint ends")
	}()
	ch := make(chan *types.Header)
	sub, err := r.client.SubscribeNewHead(ctx, ch)
	if err != nil {
		log.Panicln(err)
		return err
	}
	defer sub.Unsubscribe()
	for {
		select {
		case head := <-ch:
			// fmt.Println("header received", head.Number)
			b := BlockNumHash{
				Hash: head.Hash(),
				Num:  head.Number,
			}
			output <- b
		case <-ctx.Done():
			close(output)
			return nil
		}
	}
}
