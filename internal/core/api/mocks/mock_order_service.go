package mocks

import (
	"github.com/Owouwun/spkuznetsov/internal/core/logic/orders"
	"github.com/stretchr/testify/mock"
)

type MockOrderService struct {
	mock.Mock
}

func (m *MockOrderService) CreateNewOrder(preq *orders.PrimaryOrder) (*orders.Order, error) {
	args := m.Called(preq)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*orders.Order), args.Error(1)
}
