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
	"github.com/yktseng/portto-assignment/internal/perf"
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
	endpoints := []string{"https://data-seed-prebsc-1-s2.binance.org:8545",
		"https://data-seed-prebsc-1-s1.binance.org:8545",
		"https://data-seed-prebsc-1-s3.binance.org:8545",
		"https://data-seed-prebsc-2-s3.binance.org:8545",
	}
	rpcList := []*myeth.RPC{}

	ws := "wss://speedy-nodes-nyc.moralis.io/abd644fc46832389b55dc6d9/bsc/testnet/ws"
	for i := 0; i < len(endpoints); i++ {
		rpc := myeth.RPC{}
		result := rpc.Connect(endpoints[i])
		if !result {
			panic("failed to connect to eth endpoint")
		}
		log.Println("connected to endpoint", i)
		rpcList = append(rpcList, &rpc)
	}
	wsEndpoint := myeth.RPC{}
	result := wsEndpoint.Connect(ws)
	if !result {
		panic("failed to connect to eth endpoint")
	}
	log.Println("connected to websocket endpoint")

	var wg sync.WaitGroup
	wg.Add(*bWorkerSize + *txWorkerSize)
	ctx, cancel := context.WithCancel(context.Background())

	db := database.NewDatabase()
	err := db.Connect()
	if err != nil {
		panic(err)
	}

	bCollectors := collector.NewBlockCollector(rpcList, &wsEndpoint, db, *bWorkerSize, &wg)

	ub, err := db.GetUnfinishedBlocks(ctx)
	if err != nil {
		panic(err)
	}
	fb, err := db.GetLastRecordedBlock(ctx)
	if err != nil {
		panic(err)
	}

	mb, err := db.GetMissingBlocks(ctx)
	if err != nil {
		panic(err)
	}

	bCollectors.SetMissingBlocks(mb)
	bCollectors.SetUnfinishedBlocks(ub)
	bCollectors.SetFromBlock(fb.Add(fb, big.NewInt(1)))

	txHashes := make(chan []common.Hash, *txWorkerSize*2)
	txReceipts := make(chan []*transaction.TX, *txWorkerSize*2)
	txCollectors := collector.NewTxCollector(rpcList,
		*txWorkerSize, &wg, txHashes, txReceipts)

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Println("Ctrl-c")
		cancel()
	}()

	bPerf := make(chan int, *txWorkerSize)
	txPerf := make(chan int, *txWorkerSize)

	monitor := perf.Monitor{
		BPerf:  bPerf,
		TXPerf: txPerf,
	}
	go func() {
		monitor.Start()
	}()

	blocks := bCollectors.Start(ctx)
	txCollectors.Start(ctx)

	wg.Add(1)
	go internal.BlockHandler(ctx, &wg, db, blocks, txHashes, txReceipts)

	for i := 0; i < *txWorkerSize; i++ {
		wg.Add(1)
		go internal.TxHandler(ctx, i, &wg, db, blocks, txHashes, txReceipts, bPerf, txPerf)
	}

	wg.Wait()
	log.Println("Graceful shutdown")
}

func startPProf() {
	log.Println("start pprof")
	go http.ListenAndServe("localhost:6060", nil)
}
