package testutils

import (
	"reflect"
	"testing"
	"time"

	"github.com/seagumineko/spkuznetsov/internal/core/logic/auth"
	"github.com/seagumineko/spkuznetsov/internal/core/logic/requests"
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

type requestOption func(*requests.Request)

func WithClientName(cn string) requestOption {
	return func(r *requests.Request) {
		r.ClientName = cn
	}
}

func WithClientPhone(cp string) requestOption {
	return func(r *requests.Request) {
		r.ClientPhone = cp
	}
}

func WithAddress(a string) requestOption {
	return func(r *requests.Request) {
		r.Address = a
	}
}

func WithClientDescription(cd string) requestOption {
	return func(r *requests.Request) {
		r.ClientDescription = cd
	}
}

func WithEmployee(emp *auth.Employee) requestOption {
	return func(r *requests.Request) {
		r.Employee = emp
	}
}

func WithCancelReason(cr string) requestOption {
	return func(r *requests.Request) {
		r.CancelReason = &cr
	}
}

func WithStatus(s requests.Status) requestOption {
	return func(r *requests.Request) {
		r.Status = s
	}
}

func WithEmployeeDescription(ed string) requestOption {
	return func(r *requests.Request) {
		r.EmployeeDescription = ed
	}
}

func WithScheduledFor(sf *time.Time) requestOption {
	return func(r *requests.Request) {
		r.ScheduledFor = sf
	}
}

// Create test request with default values
func NewTestRequest(opts ...requestOption) *requests.Request {
	req := &requests.Request{
		ClientName:        "Иван Иванов",
		ClientPhone:       "+71112223344",
		Address:           "ул. Примерная, д. 1",
		ClientDescription: "Обычный клиент",
		PublicLink:        "example.com/abracadabra",
		Employee:          &auth.Employee{Name: "Петр Петров"},
		Status:            requests.StatusScheduled,
		ScheduledFor:      nil,
	}

	for _, opt := range opts {
		opt(req)
	}
	return req
}

func ValidateRequest(t *testing.T, expected, actual *requests.Request) {
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
	compare("EmployeeID", expected.Employee.ID, actual.Employee.ID)
}
