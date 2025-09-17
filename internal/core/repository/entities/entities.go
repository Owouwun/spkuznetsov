package entities

import (
	"time"

	"github.com/Owouwun/spkuznetsov/internal/core/logic/auth"
	"github.com/Owouwun/spkuznetsov/internal/core/logic/orders"
	"github.com/google/uuid"
)

type OrderEntity struct {
	ID                  uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	ClientName          string    `gorm:"not null"`
	ClientPhone         string    `gorm:"not null"`
	Address             string    `gorm:"not null"`
	ClientDescription   string
	EmployeeID          *uint
	CancelReason        string
	Status              int `gorm:"not null"`
	EmployeeDescription string
	ScheduledFor        *time.Time
	Employee            *EmployeeEntity `gorm:"foreignKey:EmployeeID;references:ID"`
}

type EmployeeEntity struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"not null"`
}

func (OrderEntity) TableName() string {
	return "public.orders"
}
func (EmployeeEntity) TableName() string {
	return "public.employees"
}

func NewOrderEntityFromLogic(ord *orders.Order) *OrderEntity {
	if ord == nil {
		return nil
	}
	oe := &OrderEntity{
		ID:                  ord.ID,
		ClientName:          ord.ClientName,
		ClientPhone:         ord.ClientPhone,
		Address:             ord.Address,
		ClientDescription:   ord.ClientDescription,
		CancelReason:        ord.CancelReason,
		Status:              int(ord.Status),
		EmployeeDescription: ord.EmployeeDescription,
		ScheduledFor:        ord.ScheduledFor,
		Employee:            (*EmployeeEntity)(ord.Employee),
	}
	if ord.Employee != nil {
		oe.EmployeeID = &ord.Employee.ID
	}
	return oe
}

func (oe *OrderEntity) ToLogicOrder() *orders.Order {
	if oe == nil {
		return nil
	}
	return &orders.Order{
		ID:                  oe.ID,
		ClientName:          oe.ClientName,
		ClientPhone:         oe.ClientPhone,
		Address:             oe.Address,
		ClientDescription:   oe.ClientDescription,
		Employee:            oe.Employee.ToLogicEmployee(),
		CancelReason:        oe.CancelReason,
		Status:              orders.Status(oe.Status),
		EmployeeDescription: oe.EmployeeDescription,
		ScheduledFor:        oe.ScheduledFor,
	}
}

func (ee *EmployeeEntity) ToLogicEmployee() *auth.Employee {
	if ee == nil {
		return nil
	}
	return &auth.Employee{
		ID:   ee.ID,
		Name: ee.Name,
	}
}
