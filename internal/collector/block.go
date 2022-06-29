package collector

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"sync"
	"time"

	"github.com/yktseng/portto-assignment/internal/block"
	"github.com/yktseng/portto-assignment/internal/myeth"
)

type BlockCollector struct {
	unfinishedBlocks []*big.Int
	fromBlock        *big.Int
	rpcList             []*myeth.RPC
	workerSize       int
	bnChan           chan *big.Int
	wg               *sync.WaitGroup
}

func NewBlockCollector(rpcList []*myeth.RPC, workerSize int, wg *sync.WaitGroup) *BlockCollector {
	return &BlockCollector{
		rpcList:        rpcList,
		workerSize: workerSize,
		bnChan:     make(chan *big.Int),
		wg:         wg,
		fromBlock:  big.NewInt(19074015),
	}
}

func (c *BlockCollector) SetUnfinishedBlocks(blocks []*big.Int) {
	c.unfinishedBlocks = blocks
}

func (c *BlockCollector) SetFromBlock(block *big.Int) {
	if block.Cmp(big.NewInt(19074015)) == 1 {
		c.fromBlock = block
	}
}

func (c *BlockCollector) Start(ctx context.Context) chan *block.Block {
	output := make(chan *block.Block, 100)
	for i := 0; i < c.workerSize; i++ {
		go c.worker(ctx, i, output)
	}
	go func() {
		// handle missing blocks
		for i := 0; i < len(c.unfinishedBlocks); i++ {
			log.Println("Unfinished block", c.unfinishedBlocks[i])
			c.bnChan <- c.unfinishedBlocks[i]
		}
		i := c.fromBlock
		log.Println("Start from block", i)
		// then start from the last block recorded in db
		for {
			num := new(big.Int).Set(i)
			select {
			case c.bnChan <- num:
				i = i.Add(i, big.NewInt(1))
			}
		}
	}()
	return output
}

func (c *BlockCollector) worker(ctx context.Context, workerNum int, output chan *block.Block) {
	defer func() {
		fmt.Println("block collector worker closed")
		c.wg.Done()
	}()
	for {
		select {
		case num := <-c.bnChan:
			var b *block.Block
			var err error
			for {
				b, err = c.rpcList[workerNum % len(c.rpcList)].GetBlock(ctx, num)
				if err != nil {
					log.Println(err)
					// if the block is not found yet, wait for 3 seconds and try again
					if err.Error() == "not found" {
						time.Sleep(3 * time.Second)
						continue
					}
					return
				}
				break
			}
		L:
			select {
			case output <- b:
			case <-ctx.Done():
				return
			default:
				time.Sleep(10 * time.Millisecond)
				break L
			}
		}
	}
}
