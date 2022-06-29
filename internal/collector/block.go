package collector

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"sync"
	"time"

	"github.com/yktseng/portto-assignment/internal/block"
	"github.com/yktseng/portto-assignment/internal/database"
	"github.com/yktseng/portto-assignment/internal/myeth"
)

type BlockCollector struct {
	unfinishedBlocks []*big.Int
	fromBlock        *big.Int
	db               *database.Database
	rpcList          []*myeth.RPC
	ws               *myeth.RPC
	workerSize       int
	bnChan           chan *big.Int
	wg               *sync.WaitGroup
	blockHeight      *big.Int
}

func NewBlockCollector(rpcList []*myeth.RPC, ws *myeth.RPC,
	db *database.Database,
	workerSize int, wg *sync.WaitGroup) *BlockCollector {
	return &BlockCollector{
		rpcList:    rpcList,
		ws:         ws,
		workerSize: workerSize,
		bnChan:     make(chan *big.Int, workerSize*2),
		db:         db,
		wg:         wg,
		fromBlock:  big.NewInt(20613750),
	}
}

func (c *BlockCollector) SetUnfinishedBlocks(blocks []*big.Int) {
	c.unfinishedBlocks = blocks
}

func (c *BlockCollector) SetFromBlock(block *big.Int) {
	if block.Cmp(big.NewInt(20613750)) == 1 {
		c.fromBlock = block
	}
}

func (c *BlockCollector) Start(ctx context.Context) chan *block.Block {
	output := make(chan *block.Block, 100)
	for i := 0; i < c.workerSize; i++ {
		go c.worker(ctx, i, output)
	}

	newBlockInfo := make(chan myeth.BlockNumHash)
	go func() {
		err := c.ws.SubscribeNewHeaders(ctx, newBlockInfo)
		if err != nil {
			log.Panicln("failed to subscribe for new headers")
		}
	}()

	blockHeight, err := c.rpcList[0].GetNewestBlock(ctx)
	if err != nil {
		log.Panicln(err)
	}
	c.blockHeight = blockHeight

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
			if num.Cmp(c.blockHeight) < 1 {
				select {
				case c.bnChan <- num:
					i = i.Add(i, big.NewInt(1))
				default:
				}
			}
			select {
			case b, ok := <-newBlockInfo:
				if !ok {
					return
				}
				log.Println("New block header", b.Num)
				c.blockHeight = b.Num
				if i.Cmp(c.blockHeight) >= 1 {
					c.bnChan <- b.Num
				}
			default:
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
		J:
		select {
		case num := <-c.bnChan:
			for {
			// log.Println("handle block", num)
			var b *block.Block
			var err error
			b, err = c.rpcList[workerNum%len(c.rpcList)].GetBlock(ctx, num)
			if err != nil {
				log.Println("block", num, err)
				// if the block is not found yet, wait for 3 seconds and try again
				if err.Error() == "not found" {
					break
				}
				return
			}
			if c.blockHeight.Cmp(big.NewInt(0)) > 0 { // see recent 20 blocks as unconfirmed
				diff := big.NewInt(0).Sub(c.blockHeight, num)
				if diff.Cmp(big.NewInt(20)) < 0 {
					fmt.Println("unconfirmed block", c.blockHeight.Int64(), num.Int64(), diff)
					b.Stable = false
				}
			}
		L:
			select {
			case output <- b:
				if !b.Stable {
					// get block number (b.Num - 20) to replace the unconfirmed one
					num = big.NewInt(0).Sub(num, big.NewInt(20))
					log.Println("fetch block", num, "again")
					continue
				}
				break J
			case <-ctx.Done():
				return
			default:
				time.Sleep(10 * time.Millisecond)
				break L
			}
		}
		case <-ctx.Done():
			return
		}
	}
}
