package requests

import (
	"time"

	"github.com/Owouwun/ipkuznetsov/internal/core/auth"
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

type PrimaryRequest struct {
	ClientName        string `json:"client_name"`
	ClientPhone       string `json:"client_phone"`
	Address           string `json:"address"`
	ClientDescription string `json:"client_description"`
}

type Request struct {
	// Immutable
	ID                int64          `json:"id" compare:"skip"`
	ClientName        string         `json:"client_name"`
	ClientPhone       string         `json:"client_phone"`
	Address           string         `json:"address"`
	ClientDescription string         `json:"client_description"`
	PublicLink        string         `json:"public_link" tests:"exception"`
	Employee          *auth.Employee `json:"employee_id"`
	CancelReason      string         `json:"cancel_reason,omitempty"`

	// Mutable
	Status              Status     `json:"status"`
	EmployeeDescription string     `json:"employee_description"`
	ScheduledFor        *time.Time `json:"scheduled_for"`
}

type RequestPatcher struct {
	ClientName          *string    `json:"client_name,omitempty"`
	ClientPhone         *string    `json:"client_phone,omitempty"`
	Address             *string    `json:"address,omitempty"`
	ClientDescription   *string    `json:"client_description,omitempty"`
	CancelReason        *string    `json:"cancel_reason,omitempty"`
	Status              Status     `json:"status,omitempty"`
	EmployeeDescription *string    `json:"employee_description,omitempty"`
	ScheduledFor        *time.Time `json:"scheduled_for,omitempty"`
	ConfirmedSchedule   *bool      `json:"confirmed_schedule,omitempty"`
}
