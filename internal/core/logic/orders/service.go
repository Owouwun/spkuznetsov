package orders

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderService interface {
	CreateOrder(ctx context.Context, pord *PrimaryOrder) (uuid.UUID, error)
	GetOrder(ctx context.Context, id uuid.UUID) (uuid.UUID, error)
}

func CreateNewOrder(db *gorm.DB, pord *PrimaryOrder) (uuid.UUID, error) {
	order, err := pord.CreateNewOrder()
	if err != nil {
		return uuid.Nil, err
	}

	result := db.Create(&order)
	if result.Error != nil {
		return uuid.Nil, result.Error
	}
	return order.ID, nil
}
