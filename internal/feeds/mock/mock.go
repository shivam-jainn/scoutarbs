package feeds

import (
	"context"

	"github.com/shivam-jainn/scoutarbs/internal/models"
)

func RunMockFeed(ctx context.Context, out chan<- models.Snapshot) {
	defer close(out)

	
}
