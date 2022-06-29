package perf

import (
	"fmt"
	"log"
	"time"
)

type Monitor struct {
	BPerf      chan int
	TXPerf     chan int
	bConsumed  int
	txConsumed int
}

func (m *Monitor) Start() {
	ticker := time.NewTicker(time.Minute)
	for {
		select {
		case <-ticker.C:
			perf := NewPerf(time.Minute, m.bConsumed, m.txConsumed)
			log.Println(perf.Summary())
			m.bConsumed = 0
			m.txConsumed = 0
		case bc := <-m.BPerf:
			m.bConsumed += bc
		case tc := <-m.TXPerf:
			m.txConsumed += tc
		default:
			time.Sleep(1 * time.Second)
		}
	}
}

type Perf struct {
	Duration      time.Duration
	BlockConsumed int
	TXConsumed    int
}

func NewPerf(duration time.Duration, bc, tc int) Perf {
	return Perf{
		Duration:      duration,
		BlockConsumed: bc,
		TXConsumed:    tc,
	}
}

func (p *Perf) Summary() string {
	avgBc := float64(p.BlockConsumed) / p.Duration.Minutes()
	avgTc := float64(p.TXConsumed) / p.Duration.Minutes()
	return fmt.Sprintf("average %f blocks and %f txs consumed per minute", avgBc, avgTc)
}
