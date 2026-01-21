package mock

import (
	"context"
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/shivam-jainn/scoutarbs/internal/models"
)

func RunMockFeed(ctx context.Context, out chan<- models.TickerData) {
	defer close(out)

	minPrice := 0.5
	maxPrice := 400000.0

	price := minPrice + rand.Float64()*(maxPrice-minPrice)
	maxDelta := 0.005 // 0.5% per tick
	seq := int64(1)

	for {
		// Base interval + jitter
		base := 300 * time.Millisecond
		jitter := time.Duration(rand.Int64N(200)-100) * time.Millisecond
		delay := base + jitter

		// Occasional lag spike (network / venue hiccup)
		if rand.Float64() < 0.02 { // 2% chance
			delay += time.Duration(rand.Int64N(2_000)) * time.Millisecond
		}

		timer := time.NewTimer(delay)

		select {
		case <-ctx.Done():
			timer.Stop()
			return

		case <-timer.C:
			timer.Stop()

			// Bursty market behavior
			burstCount := 1
			if rand.Float64() < 0.1 { // 10% chance of burst
				burstCount = int(rand.Int64N(4) + 2) // 2–5 updates
			}

			for i := 0; i < burstCount; i++ {
				delta := (rand.Float64()*2 - 1) * maxDelta
				price = price * (1 + delta)

				if price < minPrice {
					price = minPrice
				}
				if price > maxPrice {
					price = maxPrice
				}

				spread := rand.Float64()*0.5 + 0.01

				// Venue clock skew (±5ms)
				skew := time.Duration(rand.Int64N(10)-5) * time.Millisecond

				tick := models.TickerData{
					Venue:     models.VenueMock,
					Bid:       price - spread,
					Ask:       price + spread,
					BidSize:   rand.Float64() * float64(rand.Int64N(10_000_000)),
					AskSize:   rand.Float64() * float64(rand.Int64N(10_000_000)),
					Seq:       seq,
					Timestamp: time.Now().Add(skew),
				}

				seq++

				// Backpressure-safe send
				select {
				case out <- tick:
				case <-ctx.Done():
					fmt.Println("Shutting down mock feed")
					return
				}
			}
		}
	}
}
