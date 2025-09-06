package requests

import (
	"reflect"
	"testing"
	"time"

	"github.com/Owouwun/ipkuznetsov/internal/core/auth"
)

const (
	baseClientName          = "Client Name"
	baseClientPhone         = "+71112223344"
	baseAddress             = "Test Address"
	baseClientDescription   = "Test Client Description"
	baseEmployeeDescription = "Base Emp Description"
	emptyCancelReason       = ""
	filledCancelReason      = "Клиент передумал"
)

var (
	tomorrow       = time.Now().Add(24 * time.Hour)
	threeDaysLater = time.Now().Add(72 * time.Hour)
)

type requestOption func(*Request)

func withClientName(cn string) requestOption {
	return func(r *Request) {
		r.ClientName = cn
	}
}

func withClientPhone(cp string) requestOption {
	return func(r *Request) {
		r.ClientPhone = cp
	}
}

func withAddress(a string) requestOption {
	return func(r *Request) {
		r.Address = a
	}
}

func withClientDescription(cd string) requestOption {
	return func(r *Request) {
		r.ClientDescription = cd
	}
}

func withEmployee(emp *auth.Employee) requestOption {
	return func(r *Request) {
		r.Employee = emp
	}
}

func withCancelReason(cr string) requestOption {
	return func(r *Request) {
		r.CancelReason = &cr
	}
}

func withStatus(s Status) requestOption {
	return func(r *Request) {
		r.Status = s
	}
}

func withEmployeeDescription(ed string) requestOption {
	return func(r *Request) {
		r.EmployeeDescription = ed
	}
}

func withScheduledFor(sf *time.Time) requestOption {
	return func(r *Request) {
		r.ScheduledFor = sf
	}
}

// Create test request with default values
func newTestRequest(opts ...requestOption) *Request {
	req := &Request{
		ClientName:        "Иван Иванов",
		ClientPhone:       "+71112223344",
		Address:           "ул. Примерная, д. 1",
		ClientDescription: "Обычный клиент",
		PublicLink:        "example.com/abracadabra",
		Employee:          &auth.Employee{Name: "Петр Петров"},
		Status:            StatusScheduled,
		ScheduledFor:      &threeDaysLater,
	}

	for _, opt := range opts {
		opt(req)
	}
	return req
}

func assertError(t *testing.T, expected, actual error) {
	if actual != expected {
		t.Errorf("Expected error: '%v', got: '%v'", expected, actual)
	}
}

func validateRequest(t *testing.T, expected, actual *Request) {
	if (expected != nil) && (actual == nil) {
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
	compare("Employee", expected.Employee, actual.Employee)
	compare("CancelReason", expected.CancelReason, actual.CancelReason)
	compare("Status", expected.Status, actual.Status)
	compare("EmployeeDescription", expected.EmployeeDescription, actual.EmployeeDescription)
	compare("ScheduledFor", expected.ScheduledFor, actual.ScheduledFor)
}

func TestNewRequest(t *testing.T) {
	basePrimaryRequest := &PrimaryRequest{
		ClientName:        baseClientName,
		ClientPhone:       baseClientPhone,
		Address:           baseAddress,
		ClientDescription: baseClientDescription,
	}

	cases := []struct {
		name   string
		pReq   *PrimaryRequest
		expReq *Request
		expErr error
	}{
		{
			name: "Успешное создание новой заявки",
			pReq: basePrimaryRequest,
			expReq: newTestRequest(
				withClientName(basePrimaryRequest.ClientName),
				withClientPhone(basePrimaryRequest.ClientPhone),
				withAddress(basePrimaryRequest.Address),
				withClientDescription(basePrimaryRequest.ClientDescription),
				withEmployee(nil),
				withStatus(StatusNew),
				withScheduledFor(nil),
			),
			expErr: nil,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req, err := c.pReq.CreateNewRequest()
			assertError(t, c.expErr, err)
			validateRequest(t, c.expReq, req)
		})
	}
}

func TestPreschedule(t *testing.T) {
	cases := []struct {
		name   string
		req    *Request
		date   *time.Time
		expReq *Request
		expErr error
	}{
		{
			name: "Успешная попытка назначить предварительную дату новой заявки",
			req: newTestRequest(
				withStatus(StatusNew),
				withScheduledFor(nil),
			),
			date: &tomorrow,
			expReq: newTestRequest(
				withStatus(StatusPrescheduled),
				withScheduledFor(&tomorrow),
			),
			expErr: nil,
		},
		{
			name: "Успешная попытка переназначить предварительную дату новой заявки",
			req: newTestRequest(
				withStatus(StatusPrescheduled),
				withScheduledFor(&tomorrow),
			),
			date: &threeDaysLater,
			expReq: newTestRequest(
				withStatus(StatusPrescheduled),
				withScheduledFor(&threeDaysLater),
			),
			expErr: nil,
		},
		{
			name: "Успешная попытка назначить предварительную дату заявки после частичных работ",
			req: newTestRequest(
				withStatus(StatusInProgress),
				withScheduledFor(nil),
			),
			date: &tomorrow,
			expReq: newTestRequest(
				withStatus(StatusPrescheduled),
				withScheduledFor(&tomorrow),
			),
			expErr: nil,
		},
		{
			name: "Попытка назначить предварительную дату для отменённой заявки",
			req: newTestRequest(
				withStatus(StatusCanceled),
				withScheduledFor(nil),
			),
			date: &tomorrow,
			expReq: newTestRequest(
				withStatus(StatusCanceled),
				withScheduledFor(nil),
			),
			expErr: ErrActionNotPermittedByStatus,
		},
		{
			name: "Попытка назначить предварительную дату для выполненной заявки",
			req: newTestRequest(
				withStatus(StatusDone),
				withScheduledFor(nil),
			),
			date: &tomorrow,
			expReq: newTestRequest(
				withStatus(StatusDone),
				withScheduledFor(nil),
			),
			expErr: ErrActionNotPermittedByStatus,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.req.Preschedule(c.date)
			assertError(t, c.expErr, err)
			validateRequest(t, c.expReq, c.req)
		})
	}
}

func TestAssign(t *testing.T) {
	employee := &auth.Employee{
		Name: "Николай Николаев",
	}
	cases := []struct {
		name   string
		req    *Request
		emp    *auth.Employee
		expReq *Request
		expErr error
	}{
		{
			name: "Успешная попытка назначить сотрудника на новую заявку",
			req: newTestRequest(
				withEmployee(nil),
				withStatus(StatusPrescheduled),
			),
			emp: employee,
			expReq: newTestRequest(
				withEmployee(employee),
				withStatus(StatusAssigned),
			),
			expErr: nil,
		},
		{
			name: "Попытка назначить сотрудника на отменённую заявку",
			req: newTestRequest(
				withEmployee(nil),
				withStatus(StatusCanceled),
			),
			emp: employee,
			expReq: newTestRequest(
				withEmployee(nil),
				withStatus(StatusCanceled),
			),
			expErr: ErrActionNotPermittedByStatus,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.req.Assign(c.emp)
			assertError(t, c.expErr, err)
			validateRequest(t, c.expReq, c.req)
		})
	}
}

func TestSchedule(t *testing.T) {
	cases := []struct {
		name     string
		req      *Request
		schedule *time.Time
		expReq   *Request
		expErr   error
	}{
		{
			name: "Успешное планирование даты работ",
			req: newTestRequest(
				withStatus(StatusAssigned),
				withScheduledFor(&threeDaysLater),
			),
			schedule: &tomorrow,
			expReq: newTestRequest(
				withStatus(StatusScheduled),
				withScheduledFor(&tomorrow),
			),
			expErr: nil,
		},
		{
			name: "Успешное планирование даты новых работ",
			req: newTestRequest(
				withStatus(StatusInProgress),
				withScheduledFor(nil),
			),
			schedule: &tomorrow,
			expReq: newTestRequest(
				withStatus(StatusScheduled),
				withScheduledFor(&tomorrow),
			),
			expErr: nil,
		},
		{
			name: "Попытка запланировать выполненные даты работы",
			req: newTestRequest(
				withStatus(StatusDone),
			),
			schedule: &tomorrow,
			expReq: newTestRequest(
				withStatus(StatusDone),
			),
			expErr: ErrActionNotPermittedByStatus,
		},
		{
			name: "Попытка запланировать отменённые работы",
			req: newTestRequest(
				withStatus(StatusCanceled),
			),
			schedule: &threeDaysLater,
			expReq: newTestRequest(
				withStatus(StatusCanceled),
			),
			expErr: ErrActionNotPermittedByStatus,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.req.Schedule(c.schedule)
			assertError(t, c.expErr, err)
			validateRequest(t, c.expReq, c.req)
		})
	}
}

func TestConfirmSchedule(t *testing.T) {
	cases := []struct {
		name   string
		req    *Request
		expReq *Request
		expErr error
	}{
		{
			name: "Успешное подтверждение даты работ",
			req: newTestRequest(
				withStatus(StatusAssigned),
			),
			expReq: newTestRequest(
				withStatus(StatusScheduled),
			),
			expErr: nil,
		},
		{
			name: "Попытка подтвердить работы без предварительной даты",
			req: newTestRequest(
				withScheduledFor(nil),
			),
			expReq: newTestRequest(
				withScheduledFor(nil),
			),
			expErr: ErrActionNotPermittedByStatus,
		},
		{
			name: "Попытка подтвердить отменённые работы",
			req: newTestRequest(
				withStatus(StatusCanceled),
			),
			expReq: newTestRequest(
				withStatus(StatusCanceled),
			),
			expErr: ErrActionNotPermittedByStatus,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.req.ConfirmSchedule()
			assertError(t, c.expErr, err)
			validateRequest(t, c.expReq, c.req)
		})
	}
}

func TestProgress(t *testing.T) {
	cases := []struct {
		name    string
		req     *Request
		empDesc string
		expReq  *Request
		expErr  error
	}{
		{
			name: "Успешный прогресс заявки",
			req: newTestRequest(
				withStatus(StatusScheduled),
				withScheduledFor(&tomorrow),
			),
			empDesc: baseEmployeeDescription,
			expReq: newTestRequest(
				withStatus(StatusInProgress),
				withScheduledFor(nil),
				withEmployeeDescription(baseEmployeeDescription),
			),
			expErr: nil,
		},
		{
			name: "Попытка прогресса в выполненной заявке",
			req: newTestRequest(
				withStatus(StatusDone),
			),
			empDesc: "Новое описание",
			expReq: newTestRequest(
				withStatus(StatusDone),
			),
			expErr: ErrActionNotPermittedByStatus,
		},
		{
			name: "Попытка прогресса в отменённой заявке",
			req: newTestRequest(
				withStatus(StatusCanceled),
			),
			empDesc: "Новое описание",
			expReq: newTestRequest(
				withStatus(StatusCanceled),
			),
			expErr: ErrActionNotPermittedByStatus,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.req.Progress(c.empDesc)
			assertError(t, c.expErr, err)
			validateRequest(t, c.expReq, c.req)
		})
	}
}

func TestComplete(t *testing.T) {
	cases := []struct {
		name   string
		req    *Request
		expReq *Request
		expErr error
	}{
		{
			name: "Успешное завершение заявки",
			req: newTestRequest(
				withStatus(StatusInProgress),
			),
			expReq: newTestRequest(
				withStatus(StatusDone),
			),
			expErr: nil,
		},
		{
			name: "Попытка завершения новой заявки",
			req: newTestRequest(
				withStatus(StatusNew),
			),
			expReq: newTestRequest(
				withStatus(StatusNew),
			),
			expErr: ErrActionNotPermittedByStatus,
		},
		{
			name: "Попытка завершения отменённой заявки",
			req: newTestRequest(
				withStatus(StatusCanceled),
			),
			expReq: newTestRequest(
				withStatus(StatusCanceled),
			),
			expErr: ErrActionNotPermittedByStatus,
		},
		{
			name: "Попытка завершения завершённой заявки",
			req: newTestRequest(
				withStatus(StatusDone),
			),
			expReq: newTestRequest(
				withStatus(StatusDone),
			),
			expErr: ErrActionNotPermittedByStatus,
		},
		{
			name: "Попытка завершения закрытой заявки",
			req: newTestRequest(
				withStatus(StatusPaid),
			),
			expReq: newTestRequest(
				withStatus(StatusPaid),
			),
			expErr: ErrActionNotPermittedByStatus,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.req.Complete()
			assertError(t, c.expErr, err)
			validateRequest(t, c.expReq, c.req)
		})
	}
}

func TestClose(t *testing.T) {
	cases := []struct {
		name   string
		req    *Request
		expReq *Request
		expErr error
	}{
		{
			name: "Успешное закрытие выполненной заявки",
			req: newTestRequest(
				withStatus(StatusDone),
			),
			expReq: newTestRequest(
				withStatus(StatusPaid),
			),
			expErr: nil,
		},
		{
			name: "Попытка закрытия назначенной заявки",
			req: newTestRequest(
				withStatus(StatusScheduled),
			),
			expReq: newTestRequest(
				withStatus(StatusScheduled),
			),
			expErr: ErrActionNotPermittedByStatus,
		},
		{
			name: "Попытка закрытия отменённой заявки",
			req: newTestRequest(
				withStatus(StatusCanceled),
			),
			expReq: newTestRequest(
				withStatus(StatusCanceled),
			),
			expErr: ErrActionNotPermittedByStatus,
		},
		{
			name: "Попытка закрытия закрытой заявки",
			req: newTestRequest(
				withStatus(StatusPaid),
			),
			expReq: newTestRequest(
				withStatus(StatusPaid),
			),
			expErr: ErrActionNotPermittedByStatus,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.req.Close()
			assertError(t, c.expErr, err)
			validateRequest(t, c.expReq, c.req)
		})
	}
}

func TestCancel(t *testing.T) {
	cases := []struct {
		name         string
		req          *Request
		cancelReason string
		expReq       *Request
		expErr       error
	}{
		{
			name: "Успешная отмена новой заявки",
			req: newTestRequest(
				withStatus(StatusNew),
			),
			cancelReason: filledCancelReason,
			expReq: newTestRequest(
				withStatus(StatusCanceled),
				withCancelReason(filledCancelReason),
				withScheduledFor(nil),
			),
			expErr: nil,
		},
		{
			name: "Успешная отмена запланированной заявки",
			req: newTestRequest(
				withStatus(StatusScheduled),
			),
			cancelReason: filledCancelReason,
			expReq: newTestRequest(
				withStatus(StatusCanceled),
				withCancelReason(filledCancelReason),
				withScheduledFor(nil),
			),
			expErr: nil,
		},
		{
			name: "Попытка отмены отменённой заявки",
			req: newTestRequest(
				withStatus(StatusCanceled),
				withCancelReason(filledCancelReason),
			),
			cancelReason: "Другая причина отмены",
			expReq: newTestRequest(
				withStatus(StatusCanceled),
				withCancelReason(filledCancelReason),
			),
			expErr: ErrActionNotPermittedByStatus,
		},
		{
			name: "Попытка отмены оплаченной заявки",
			req: newTestRequest(
				withStatus(StatusPaid),
			),
			cancelReason: filledCancelReason,
			expReq: newTestRequest(
				withStatus(StatusPaid),
			),
			expErr: ErrActionNotPermittedByStatus,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.req.Cancel(c.cancelReason)
			assertError(t, c.expErr, err)
			validateRequest(t, c.expReq, c.req)
		})
	}
}

func TestPatch(t *testing.T) {
	patchedClientName := "Patched Cliend Name"
	patchedClientPhone := "+72222222222"
	patchedAddress := "Patched Test Address"
	patchedCliendDescription := "Patched Cliend Descr"
	patchedEmployeeDescription := "Patched Emp Descr"

	cases := []struct {
		name          string
		req           *Request
		patchedFields *RequestPatcher
		expReq        *Request
		expErr        error
	}{
		{
			name: "Успешная модификация полей заявки",
			req:  newTestRequest(),
			patchedFields: &RequestPatcher{
				ClientName:          &patchedClientName,
				ClientPhone:         &patchedClientPhone,
				Address:             &patchedAddress,
				ClientDescription:   &patchedCliendDescription,
				EmployeeDescription: &patchedEmployeeDescription,
			},
			expReq: newTestRequest(
				withClientName(patchedClientName),
				withClientPhone(patchedClientPhone),
				withAddress(patchedAddress),
				withClientDescription(patchedCliendDescription),
				withEmployeeDescription(patchedEmployeeDescription),
			),
			expErr: nil,
		},
		{
			name: "Попытка модификации отменённой заявки",
			req: newTestRequest(
				withStatus(StatusCanceled),
			),
			patchedFields: &RequestPatcher{
				ClientName:          &patchedClientName,
				ClientPhone:         &patchedClientPhone,
				Address:             &patchedAddress,
				ClientDescription:   &patchedCliendDescription,
				EmployeeDescription: &patchedEmployeeDescription,
			},
			expReq: newTestRequest(
				withStatus(StatusCanceled),
			),
			expErr: ErrActionNotPermittedByStatus,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.req.Patch(c.patchedFields)
			assertError(t, c.expErr, err)
			validateRequest(t, c.expReq, c.req)
		})
	}
}
