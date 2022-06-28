package main

import (
	"context"
	"flag"
	"log"
	"math/big"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/ethereum/go-ethereum/common"
	"github.com/yktseng/portto-assignment/internal"
	"github.com/yktseng/portto-assignment/internal/collector"
	"github.com/yktseng/portto-assignment/internal/database"
	"github.com/yktseng/portto-assignment/internal/myeth"
	transaction "github.com/yktseng/portto-assignment/internal/tx"
)

var pprof = flag.Bool("pprof", false, "enable pprof")
var bWorkerSize = flag.Int("block-workers", 1, "number of block collectors")
var txWorkerSize = flag.Int("tx-workers", 4, "number of tx collectors")



func main() {

	flag.Parse()
	if *pprof {
		startPProf()
	}

	rpc := myeth.RPC{}
	result := rpc.Connect("https://data-seed-prebsc-1-s1.binance.org:8545")
	if !result {
		panic("failed to connect to eth endpoint")
	}
	log.Println("connected to endpoint")

	var wg sync.WaitGroup
	wg.Add(*bWorkerSize + *txWorkerSize)
	ctx, cancel := context.WithCancel(context.Background())

	db := database.NewDatabase()
	err := db.Connect()
	if err != nil {
		panic(err)
	}

	bCollectors := collector.NewBlockCollector(&rpc, *bWorkerSize, &wg)
	
	ub, err := db.GetUnfinishedBlocks(ctx)
	if err != nil {
		panic(err)
	}
	fb, err := db.GetLastRecordedBlock(ctx)
	if err != nil {
		panic(err)
	}
	bCollectors.SetUnfinishedBlocks(ub)
	bCollectors.SetFromBlock(fb.Add(fb, big.NewInt(1)))

	txHashes := make(chan common.Hash, 48)
	txReceipts := make(chan *transaction.TX, 100)
	txCollectors := collector.NewTxCollector(&rpc,
		*txWorkerSize, &wg, txHashes, txReceipts)

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Println("Ctrl-c")
		cancel()
	}()
	blocks := bCollectors.Start(ctx)
	txCollectors.Start(ctx)

	wg.Add(1)
	go internal.BlockHandler(ctx, &wg, db, blocks, txHashes, txReceipts)

	for i := 0; i < *txWorkerSize * 3; i++ {
		wg.Add(1)
		go internal.TxHandler(ctx, &wg, db, blocks, txHashes, txReceipts)
	}

	wg.Wait()
	log.Println("Graceful shutdown")
}

func startPProf() {
	log.Println("start pprof")
	go http.ListenAndServe("localhost:6060", nil)
}
