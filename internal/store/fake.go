package store

import (
	"context"
	"github.com/stretchr/testify/mock"
)

type Fake struct {
	mock.Mock
}

func NewFake() *Fake {
	return &Fake{}
}

func (s *Fake) Save(ctx context.Context, purchase StockRequest) error {
	arguments := s.Called(ctx, purchase)
	return arguments.Error(0)
}

func (s *Fake) StubSave() *mock.Call {
	return s.On("Save", mock.Anything, mock.Anything)
}
