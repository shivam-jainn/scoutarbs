package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/shivam-jainn/scoutarbs/internal/feeds"
	"golang.org/x/sync/errgroup"
)

func main() {
	fmt.Println("an arbitrage trading scouting engine")
	fmt.Println("starting ...")

	rootCtx, rootCancel := context.WithCancel(context.Background())
	defer rootCancel()

	g, ctx := errgroup.WithContext(rootCtx)

	// run feeds under errgroup so errors cancel context and main can wait
	g.Go(func() error { return feeds.Run(ctx) })

	// wait or handle graceful shutdown on signal
	errCh := make(chan error, 1)
	go func() { errCh <- g.Wait() }()

	select {
	case err := <-errCh:
		if err != nil {
			log.Printf("fatal: %v", err)
		}
	case <-ctx.Done():
		log.Println("shutdown signal received — waiting for graceful shutdown")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		done := make(chan struct{})
		go func() { _ = g.Wait(); close(done) }()

		select {
		case <-done:
			log.Println("graceful shutdown complete")
		case <-shutdownCtx.Done():
			log.Println("graceful shutdown timed out; exiting")
		}
	}

}
