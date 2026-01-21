package feeds

import (
	"context"
	"fmt"
	"sync"

	"github.com/shivam-jainn/scoutarbs/internal/feeds/mock"
	"github.com/shivam-jainn/scoutarbs/internal/models"
)

func Run(ctx context.Context) error {
	out := make(chan models.TickerData)
	var wg sync.WaitGroup
	producers := 2

	wg.Add(producers)
	for i := 0; i < producers; i++ {
		go func() {
			defer wg.Done()
			mock.RunMockFeed(ctx, out)
		}()
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case update, ok := <-out:
			if !ok {
				return nil
			}
			fmt.Println(update)
		}
	}
}
