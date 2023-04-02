package store

import "context"

type Interface interface {
	Save(ctx context.Context, purchase StockRequest) error
}
