package testutils

import (
	"reflect"
	"testing"
	"time"

	"github.com/Owouwun/spkuznetsov/internal/core/logic/auth"
	"github.com/Owouwun/spkuznetsov/internal/core/logic/orders"
)

const (
	ClientName          = "Client Name"             // Default value
	ClientPhone         = "+71112223344"            // Default value
	Address             = "Test Address"            // Default value
	ClientDescription   = "Test Client Description" // Default value
	EmployeeDescription = "Base Emp Description"    // Default value
	EmptyCancelReason   = ""                        // Default value
	FilledCancelReason  = "Клиент передумал"
)

func GetNDaysLater(n int) time.Time {
	return time.Now().Add(time.Duration(24*n) * time.Hour)
}

type OrderOption func(*orders.Order)

func WithClientName(cn string) OrderOption {
	return func(r *orders.Order) {
		r.ClientName = cn
	}
}

func WithClientPhone(cp string) OrderOption {
	return func(r *orders.Order) {
		r.ClientPhone = cp
	}
}

func WithAddress(a string) OrderOption {
	return func(r *orders.Order) {
		r.Address = a
	}
}

func WithClientDescription(cd string) OrderOption {
	return func(r *orders.Order) {
		r.ClientDescription = cd
	}
}

func WithEmployee(emp *auth.Employee) OrderOption {
	return func(r *orders.Order) {
		r.Employee = emp
	}
}

func WithCancelReason(cr string) OrderOption {
	return func(r *orders.Order) {
		r.CancelReason = cr
	}
}

func WithStatus(s orders.Status) OrderOption {
	return func(r *orders.Order) {
		r.Status = s
	}
}

func WithEmployeeDescription(ed string) OrderOption {
	return func(r *orders.Order) {
		r.EmployeeDescription = ed
	}
}

func WithScheduledFor(sf *time.Time) OrderOption {
	return func(r *orders.Order) {
		r.ScheduledFor = sf
	}
}

// Create test Order with default values
func NewTestOrder(opts ...OrderOption) *orders.Order {
	req := &orders.Order{
		ClientName:        "Иван Иванов",
		ClientPhone:       "+71112223344",
		Address:           "ул. Примерная, д. 1",
		ClientDescription: "Обычный клиент",
		Employee:          &auth.Employee{ID: 1, Name: "Петр Петров"},
		Status:            orders.StatusScheduled,
		ScheduledFor:      nil,
	}

	for _, opt := range opts {
		opt(req)
	}
	return req
}

func ValidateOrder(t *testing.T, expected, actual *orders.Order) {
	if expected == nil {
		if actual == nil {
			return
		}
		t.Errorf("Got non-nil value, while expected nil")
		return
	}
	if actual == nil {
		t.Errorf("Got nil value, while expected non-nil")
		return
	}

	compare := func(field string, a, b interface{}) {
		if !reflect.DeepEqual(a, b) {
			t.Errorf("Field '%s': expected '%v', got '%v'", field, a, b)
		}
	}

	compare("ClientName", expected.ClientName, actual.ClientName)
	compare("ClientPhone", expected.ClientPhone, actual.ClientPhone)
	compare("Address", expected.Address, actual.Address)
	compare("ClientDescription", expected.ClientDescription, actual.ClientDescription)
	compare("CancelReason", expected.CancelReason, actual.CancelReason)
	compare("Status", expected.Status, actual.Status)
	compare("EmployeeDescription", expected.EmployeeDescription, actual.EmployeeDescription)
	compare("ScheduledFor", expected.ScheduledFor, actual.ScheduledFor)

	if expected.Employee == nil {
		if actual.Employee == nil {
			return
		}
		t.Errorf("Got non-nil value, while expected nil")
		return
	}
	if actual.Employee == nil {
		t.Errorf("Got nil value, while expected non-nil")
		return
	}
	compare("EmployeeName", expected.Employee.Name, actual.Employee.Name)
}
