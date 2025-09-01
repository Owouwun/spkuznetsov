package requests

import (
	"time"
)

type RequestStatus string

const (
	StatusNew          RequestStatus = "Оформлена"
	StatusPrescheduled RequestStatus = "Назначена предварительная дата работ"
	StatusAssigned     RequestStatus = "Назначен сотрудник"
	StatusScheduled    RequestStatus = "Назначены работы"
	StatusInProgress   RequestStatus = "Работы частично проведены"
	StatusDone         RequestStatus = "Выполнена"
	StatusPaid         RequestStatus = "Оплачена"
	StatusCanceled     RequestStatus = "Отменена"
)

type PrimaryRequest struct {
	ClientName        string `json:"client_name"`
	ClientPhone       string `json:"client_phone"`
	Address           string `json:"address"`
	ClientDescription string `json:"client_description"`
}

type Request struct {
	// Immutable
	ID                int64  `json:"id"`
	ClientName        string `json:"client_name"`
	ClientPhone       string `json:"client_phone"`
	Address           string `json:"address"`
	ClientDescription string `json:"client_description"`
	PublicLink        string `json:"public_link"`
	EmployeeID        int64  `json:"employee_id"`
	CancelReason      string `json:"cancel_reason,omitempty"`

	// Mutable
	Status              RequestStatus `json:"status"`
	EmployeeDescription string        `json:"employee_description"`
	ScheduledFor        *time.Time    `json:"scheduled_for"`
	ConfirmedSchedule   bool          `json:"confirmed_schedule"`
	Done                bool          `json:"done"`
	Paid                bool          `json:"Paid"`
}

type RequestHistoryChange struct {
	Status              *RequestStatus `json:"status,omitempty"`
	EmployeeDescription *string        `json:"employee_description,omitempty"`
	ScheduledFor        *time.Time     `json:"scheduled_for,omitempty"`
	ConfirmedSchedule   *bool          `json:"confirmed_schedule,omitempty"`
	Done                *bool          `json:"done,omitempty"`
	Paid                *bool          `json:"paid,omitempty"`
}

type RequestHistoryEntry struct {
	HistoryID   int64                 `json:"history_id"`
	RequestID   int64                 `json:"request_id"`
	Timestamp   time.Time             `json:"timestamp"`
	ChangedByID int64                 `json:"changed_by_id"`
	Changes     *RequestHistoryChange `json:"changes"`
}
