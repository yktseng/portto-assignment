package internal

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/yktseng/portto-assignment/internal/block"
	"github.com/yktseng/portto-assignment/internal/database"
	transaction "github.com/yktseng/portto-assignment/internal/tx"
)

func newTxEntries(b *block.Block) *sync.Map {
	var v sync.Map
	for _, tx := range b.Transactions() {
		v.Store(tx, true)
	}
	return &v
}

func isEmpty(m *sync.Map) bool {
	count := 0
	m.Range(func(any, any) bool {
		count++
		return false
	})
	return count == 0
}

var blockTXMap sync.Map

func BlockHandler(ctx context.Context, wg *sync.WaitGroup, db *database.Database, blocks chan *block.Block,
	txHashes chan common.Hash, txReceipts chan *transaction.TX) {
	defer func() {
		wg.Done()
		recover()
	}()

	for {
		select {
		case block := <-blocks:
			// log.Println("received block", block.Num.Int, block.Hash)
			// create an entry in blockTXMap
			e := newTxEntries(block)
			if len(block.Transactions()) > 0 {
				blockTXMap.Store(block.Hash, e)
			}
			// log.Println("save block", block.Num.Int)
			err := db.SaveBlock(ctx, block)
			if err != nil {
				log.Panicln("failed to write block to db")
			}
			h := common.BytesToHash(common.FromHex(block.Hash))
			db.SetBlockDone(ctx, h)
			txs := block.Transactions()
			for i := 0; i < len(txs); {
				// Send TXs to tx workers
				select {
				case txHashes <- txs[i]:
					// fmt.Println(i)
					i++
				case <-ctx.Done():
					log.Println("block handler closed")
					return
				default:
					time.Sleep(10 * time.Millisecond)
				}
			}
		case <-ctx.Done():
			log.Println("block handler closed")
			return
		}
	}
}

func TxHandler(ctx context.Context, wg *sync.WaitGroup, db *database.Database, blocks chan *block.Block,
	txHashes chan common.Hash, txReceipts chan *transaction.TX) {
	defer func() {
		wg.Done()
		recover()
	}()

	for {
		select {
		case tx := <-txReceipts:
			// log.Println("receive tx", tx.Hash)
			// write tx and logs to db
			err := db.SaveTxs(ctx, []*transaction.TX{tx})
			if err != nil {
				log.Panicln("failed to write txs to db")
			}
			err = db.SaveLogs(ctx, tx.Logs)
			if err != nil {
				log.Panicln("failed to write logs to db")
			}
			// remove tx entries in blockTXMap
			h := common.BytesToHash(common.FromHex(tx.BlockHash))
			d, ok := blockTXMap.Load(tx.BlockHash)
			if !ok {
				// should not happen
				// fmt.Println(blockTXMap, h)
				log.Panicln("BUG: tx is not found in blockTXMap")
			}
			e, ok := d.(*sync.Map)
			if !ok {
				log.Panicln("BUG: tx entry has deleted")
			}
			e.Delete(common.BytesToHash(common.FromHex(tx.Hash)))
			// remove block entry if all txs are collected
			if isEmpty(e) {
				log.Println("is done: block", tx.BlockHash)
				blockTXMap.Delete(tx.BlockHash)
				db.SetBlockDone(ctx, h)
			}
		case <-ctx.Done():
			log.Println("tx handler closed")
			return
		}
	}
}
