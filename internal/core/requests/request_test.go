package requests

import (
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

func withPublicLink(pl string) requestOption {
	return func(r *Request) {
		r.PublicLink = pl
	}
}

func withEmployee(emp *auth.Employee) requestOption {
	return func(r *Request) {
		r.Employee = emp
	}
}

func withCancelReason(cr string) requestOption {
	return func(r *Request) {
		r.CancelReason = cr
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

func newRequest(opts ...requestOption) *Request {
	threeDaysLater := time.Now().Add(72 * time.Hour)
	req := &Request{
		ClientName:          "Иван Иванов",
		ClientPhone:         "+71112223344",
		Address:             "ул. Примерная, д. 1",
		ClientDescription:   "Обычный клиент",
		PublicLink:          "example.com/abracadabra",
		Employee:            &auth.Employee{Name: "Петр Петров"},
		CancelReason:        "",
		Status:              StatusScheduled,
		EmployeeDescription: "Обычное описание",
		ScheduledFor:        &threeDaysLater,
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
	if actual != expected {
		t.Errorf("Actual and expected requests don't match")
	}
}

func TestNewRequest(t *testing.T) {
	t.Error(ErrNotImplemented)
}

func TestPreschedule(t *testing.T) {
	t.Error(ErrNotImplemented)
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
			req: newRequest(
				withEmployee(nil),
				withStatus(StatusPrescheduled),
			),
			emp: employee,
			expReq: newRequest(
				withEmployee(employee),
				withStatus(StatusAssigned),
			),
			expErr: nil,
		},
		{
			name: "Попытка назначить сотрудника на отменённую заявку",
			req: newRequest(
				withEmployee(nil),
				withStatus(StatusCanceled),
			),
			emp: employee,
			expReq: newRequest(
				withEmployee(nil),
				withStatus(StatusCanceled),
			),
			expErr: ErrActionNotPermittedByStatus,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.req.Schedule(&threeDaysLater)
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
			req: newRequest(
				withStatus(StatusAssigned),
				withScheduledFor(&threeDaysLater),
			),
			schedule: &tomorrow,
			expReq: newRequest(
				withStatus(StatusScheduled),
				withScheduledFor(&tomorrow),
			),
			expErr: nil,
		},
		{
			name: "Успешное планирование даты новых работ",
			req: newRequest(
				withStatus(StatusInProgress),
				withScheduledFor(nil),
			),
			schedule: &tomorrow,
			expReq: newRequest(
				withStatus(StatusScheduled),
				withScheduledFor(&tomorrow),
			),
			expErr: nil,
		},
		{
			name: "Попытка запланировать выполненные даты работы",
			req: newRequest(
				withStatus(StatusDone),
			),
			schedule: &tomorrow,
			expReq: newRequest(
				withStatus(StatusDone),
			),
			expErr: ErrActionNotPermittedByStatus,
		},
		{
			name: "Попытка запланировать отменённые работы",
			req: newRequest(
				withStatus(StatusCanceled),
			),
			schedule: &threeDaysLater,
			expReq: newRequest(
				withStatus(StatusCanceled),
			),
			expErr: ErrActionNotPermittedByStatus,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.req.Schedule(&threeDaysLater)
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
			req: newRequest(
				withStatus(StatusAssigned),
			),
			expReq: newRequest(
				withStatus(StatusScheduled),
			),
			expErr: nil,
		},
		{
			name: "Попытка подтвердить работы без предварительной даты",
			req: newRequest(
				withScheduledFor(nil),
			),
			expReq: newRequest(
				withScheduledFor(nil),
			),
			expErr: ErrActionNotPermittedByStatus,
		},
		{
			name: "Попытка подтвердить отменённые работы",
			req: newRequest(
				withStatus(StatusCanceled),
			),
			expReq: newRequest(
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
			req: newRequest(
				withStatus(StatusScheduled),
				withScheduledFor(&threeDaysLater),
			),
			empDesc: baseEmployeeDescription,
			expReq: newRequest(
				withStatus(StatusInProgress),
				withScheduledFor(nil),
			),
			expErr: nil,
		},
		{
			name: "Попытка прогресса в выполненной заявке",
			req: newRequest(
				withStatus(StatusDone),
			),
			empDesc: "Новое описание",
			expReq: newRequest(
				withStatus(StatusDone),
			),
			expErr: ErrActionNotPermittedByStatus,
		},
		{
			name: "Попытка прогресса в отменённой заявке",
			req: newRequest(
				withStatus(StatusCanceled),
			),
			empDesc: "Новое описание",
			expReq: newRequest(
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
			req: newRequest(
				withStatus(StatusInProgress),
			),
			expReq: newRequest(
				withStatus(StatusDone),
			),
			expErr: nil,
		},
		{
			name: "Попытка завершения новой заявки",
			req: newRequest(
				withStatus(StatusNew),
			),
			expReq: newRequest(
				withStatus(StatusNew),
			),
			expErr: ErrActionNotPermittedByStatus,
		},
		{
			name: "Попытка завершения отменённой заявки",
			req: newRequest(
				withStatus(StatusCanceled),
			),
			expReq: newRequest(
				withStatus(StatusCanceled),
			),
			expErr: ErrActionNotPermittedByStatus,
		},
		{
			name: "Попытка завершения завершённой заявки",
			req: newRequest(
				withStatus(StatusDone),
			),
			expReq: newRequest(
				withStatus(StatusDone),
			),
			expErr: ErrActionNotPermittedByStatus,
		},
		{
			name: "Попытка завершения закрытой заявки",
			req: newRequest(
				withStatus(StatusPaid),
			),
			expReq: newRequest(
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
			req: newRequest(
				withStatus(StatusDone),
			),
			expReq: newRequest(
				withStatus(StatusPaid),
			),
			expErr: nil,
		},
		{
			name: "Попытка закрытия назначенной заявки",
			req: newRequest(
				withStatus(StatusScheduled),
			),
			expReq: newRequest(
				withStatus(StatusScheduled),
			),
			expErr: ErrActionNotPermittedByStatus,
		},
		{
			name: "Попытка закрытия отменённой заявки",
			req: newRequest(
				withStatus(StatusCanceled),
			),
			expReq: newRequest(
				withStatus(StatusCanceled),
			),
			expErr: ErrActionNotPermittedByStatus,
		},
		{
			name: "Попытка закрытия закрытой заявки",
			req: newRequest(
				withStatus(StatusPaid),
			),
			expReq: newRequest(
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
			req: newRequest(
				withStatus(StatusNew),
			),
			cancelReason: filledCancelReason,
			expReq: newRequest(
				withStatus(StatusNew),
				withCancelReason(filledCancelReason),
			),
			expErr: nil,
		},
		{
			name: "Успешная отмена запланированной заявки",
			req: newRequest(
				withStatus(StatusScheduled),
			),
			cancelReason: filledCancelReason,
			expReq: newRequest(
				withStatus(StatusCanceled),
				withCancelReason(filledCancelReason),
				withScheduledFor(nil),
			),
			expErr: nil,
		},
		{
			name: "Попытка отмены отменённой заявки",
			req: newRequest(
				withStatus(StatusCanceled),
				withCancelReason(filledCancelReason),
			),
			cancelReason: "Другая причина отмены",
			expReq: newRequest(
				withStatus(StatusCanceled),
				withCancelReason(filledCancelReason),
			),
			expErr: ErrActionNotPermittedByStatus,
		},
		{
			name: "Попытка отмены оплаченной заявки",
			req: newRequest(
				withStatus(StatusPaid),
			),
			cancelReason: filledCancelReason,
			expReq: newRequest(
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
			req:  newRequest(),
			patchedFields: &RequestPatcher{
				ClientName:          &patchedClientName,
				ClientPhone:         &patchedClientPhone,
				Address:             &patchedAddress,
				ClientDescription:   &patchedCliendDescription,
				EmployeeDescription: &patchedEmployeeDescription,
			},
			expReq: newRequest(
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
			req: newRequest(
				withStatus(StatusCanceled),
			),
			patchedFields: &RequestPatcher{
				ClientName:          &patchedClientName,
				ClientPhone:         &patchedClientPhone,
				Address:             &patchedAddress,
				ClientDescription:   &patchedCliendDescription,
				EmployeeDescription: &patchedEmployeeDescription,
			},
			expReq: newRequest(
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
