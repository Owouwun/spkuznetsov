package orders

import (
	"time"

	"github.com/Owouwun/spkuznetsov/internal/core/logic/auth"
	"gorm.io/gorm"
)

type Status int8

const (
	StatusNew          Status = 0  // "Оформлена"
	StatusPrescheduled Status = 1  // "Назначена предварительная дата работ"
	StatusAssigned     Status = 2  // "Назначен сотрудник"
	StatusScheduled    Status = 3  // "Назначены работы"
	StatusInProgress   Status = 4  // "Работы частично проведены"
	StatusDone         Status = 5  // "Выполнена"
	StatusPaid         Status = 6  // "Оплачена"
	StatusCanceled     Status = -1 // "Отменена"
)

func (s Status) ToString() string {
	switch s {
	case StatusNew:
		return "New"
	case StatusPrescheduled:
		return "Prescheduled"
	case StatusAssigned:
		return "Assigned"
	case StatusScheduled:
		return "Scheduled"
	case StatusInProgress:
		return "InProgress"
	case StatusDone:
		return "Done"
	case StatusPaid:
		return "Paid"
	case StatusCanceled:
		return "Canceled"
	}
	return ""
}

type PrimaryOrder struct {
	ClientName        string `json:"client_name"`
	ClientPhone       string `json:"client_phone"`
	Address           string `json:"address"`
	ClientDescription string `json:"client_description"`
}

type Order struct {
	// Immutable
	ClientName        string         `json:"client_name"`
	ClientPhone       string         `json:"client_phone"`
	Address           string         `json:"address"`
	ClientDescription string         `json:"client_description"`
	PublicLink        string         `json:"public_link"`
	Employee          *auth.Employee `gorm:"foreignKey:EmployeeID"`
	CancelReason      *string        `json:"cancel_reason,omitempty"`

	// Mutable
	Status              Status     `json:"status"`
	EmployeeDescription string     `json:"employee_description"`
	ScheduledFor        *time.Time `json:"scheduled_for"`

	// Gorm
	gorm.Model
	EmployeeID *uint64 `json:"employee_id,omitempty"`
}

type OrderPatcher struct {
	ClientName          *string `json:"client_name,omitempty"`
	ClientPhone         *string `json:"client_phone,omitempty"`
	Address             *string `json:"address,omitempty"`
	ClientDescription   *string `json:"client_description,omitempty"`
	EmployeeDescription *string `json:"employee_description,omitempty"`
}
