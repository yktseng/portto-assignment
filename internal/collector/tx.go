package collector

import (
	"context"
	"log"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/yktseng/portto-assignment/internal/myeth"
	transaction "github.com/yktseng/portto-assignment/internal/tx"
)

type TxCollector struct {
	rpc          *myeth.RPC
	workerSize   int
	txHashChan   chan common.Hash
	txResultChan chan *transaction.TX
	wg           *sync.WaitGroup
}

func NewTxCollector(rpc *myeth.RPC,
	workerSize int, wg *sync.WaitGroup,
	input chan common.Hash, output chan *transaction.TX) *TxCollector {
	return &TxCollector{
		rpc:          rpc,
		workerSize:   workerSize,
		txHashChan:   input,
		txResultChan: output,
		wg:           wg,
	}
}

func (c *TxCollector) Start(ctx context.Context) {
	for i := 0; i < c.workerSize; i++ {
		go c.worker(ctx, i)
	}
}

func (c *TxCollector) worker(ctx context.Context, workerNum int) {
	defer func() {
		c.wg.Done()
		recover()
	}()
	for {
		select {
		case hash := <-c.txHashChan:
			// fmt.Println("start handling", hash)
			t, err := c.rpc.GetTx(ctx, hash)
			if err != nil {
				log.Println(err)
				return
			}
			receipt, err := c.rpc.GetTxReceipt(ctx, hash)
			if err != nil {
				log.Println(err)
				return
			}
			tx, err := transaction.FromGethTX(t, receipt)
			if err != nil {
				log.Println(err)
				return
			}
			c.txResultChan <- tx
		case <-ctx.Done():
			log.Println("tx collector worker closed")
			return
		}
	}
}
