package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/yktseng/portto-assignment/internal/database"
	"github.com/yktseng/portto-assignment/internal/web"
)

var pprof = flag.Bool("pprof", false, "enable pprof")

func main() {

	flag.Parse()
	if *pprof {
		startPProf()
	}
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())

	db := database.NewDatabase()
	err := db.Connect()
	if err != nil {
		panic(err)
	}

	wg.Add(1)
	s := web.NewServer(ctx, db, &wg)

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Println("Ctrl-c")
		s.Stop(ctx)
		cancel()
	}()
	go func() {
		if err := s.Start(8080); err != http.ErrServerClosed {
			log.Panicln(err)
		}
	}()

	wg.Wait()
}

func startPProf() {
	log.Println("start pprof")
	go http.ListenAndServe("localhost:6060", nil)
}
