package requests

import (
	"fmt"
	"reflect"
	"testing"
	"time"
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
	threeDaysLater = time.Now().Add(72 * time.Hour)
)

// Сравнение всех полей с особой обработкой для полей с тегом compare:"skip"
func (got Request) validate(exp Request) error {
	// TODO Придумать модификацию алгоритма, чтобы он сам выбирал с чем сравнивать поля с compare:"skip"
	if got.ID == 0 {
		return fmt.Errorf("ID was not defined")
	}
	if got.PublicLink == "" {
		return fmt.Errorf("Public link was not defined")
	}
	if got.EmployeeID != nil {
		return fmt.Errorf("Employee ID must be nil")
	}

	rv := reflect.ValueOf(got)
	for i := 0; i < rv.NumField(); i++ {
		if rv.Type().Field(i).Tag.Get("compare") == "skip" {
			continue
		}
		g, e := rv.Field(i), reflect.ValueOf(exp).Field(i)
		if !reflect.DeepEqual(g.Interface(), e.Interface()) {
			return fmt.Errorf("%s: expected %v, got %v",
				rv.Type().Field(i).Name, e.Interface(), g.Interface())
		}
	}
	return nil
}

func assertError(t *testing.T, expected, actual error) {
	if actual != expected {
		t.Errorf("Expected error: '%v', got: '%v'", expected, actual)
	}
}

func validateRequest(t *testing.T, expected, actual *Request) {
	if err := actual.validate(*expected); err != nil {
		t.Error(err)
	}
}

func TestNewRequest(t *testing.T) {
	t.Error(ErrNotImplemented)
}

func TestPreschedule(t *testing.T) {
	t.Error(ErrNotImplemented)
}

func TestAssign(t *testing.T) {
	t.Error(ErrNotImplemented)
}

func TestSchedule(t *testing.T) {
	t.Error(ErrNotImplemented)
}

func TestConfirmSchedule(t *testing.T) {
	t.Error(ErrNotImplemented)
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
			req: &Request{
				ClientName:          baseClientName,
				ClientPhone:         baseClientPhone,
				Address:             baseAddress,
				ClientDescription:   baseClientDescription,
				EmployeeDescription: "",
				CancelReason:        emptyCancelReason,
				Status:              StatusScheduled,
				ScheduledFor:        &threeDaysLater,
			},
			empDesc: baseEmployeeDescription,
			expReq: &Request{
				ClientName:          baseClientName,
				ClientPhone:         baseClientPhone,
				Address:             baseAddress,
				ClientDescription:   baseClientDescription,
				EmployeeDescription: baseEmployeeDescription,
				CancelReason:        emptyCancelReason,
				Status:              StatusInProgress,
				ScheduledFor:        nil,
			},
			expErr: nil,
		},
		{
			name: "Попытка прогресса в выполненной заявке",
			req: &Request{
				ClientName:          baseClientName,
				ClientPhone:         baseClientPhone,
				Address:             baseAddress,
				ClientDescription:   baseClientDescription,
				EmployeeDescription: "",
				CancelReason:        emptyCancelReason,
				Status:              StatusDone,
				ScheduledFor:        nil,
			},
			empDesc: "Новое описание",
			expReq: &Request{
				ClientName:          baseClientName,
				ClientPhone:         baseClientPhone,
				Address:             baseAddress,
				ClientDescription:   baseClientDescription,
				EmployeeDescription: baseEmployeeDescription,
				CancelReason:        emptyCancelReason,
				Status:              StatusDone,
				ScheduledFor:        nil,
			},
			expErr: ErrActionNotPermittedByStatus,
		},
		{
			name: "Попытка прогресса в отменённой заявке",
			req: &Request{
				ClientName:          baseClientName,
				ClientPhone:         baseClientPhone,
				Address:             baseAddress,
				ClientDescription:   baseClientDescription,
				EmployeeDescription: "",
				CancelReason:        filledCancelReason,
				Status:              StatusCanceled,
				ScheduledFor:        nil,
			},
			empDesc: "Новое описание",
			expReq: &Request{
				ClientName:          baseClientName,
				ClientPhone:         baseClientPhone,
				Address:             baseAddress,
				ClientDescription:   baseClientDescription,
				EmployeeDescription: "",
				CancelReason:        filledCancelReason,
				Status:              StatusCanceled,
				ScheduledFor:        nil,
			},
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
			req: &Request{
				ClientName:          baseClientName,
				ClientPhone:         baseClientPhone,
				Address:             baseAddress,
				ClientDescription:   baseClientDescription,
				EmployeeDescription: baseEmployeeDescription,
				CancelReason:        emptyCancelReason,
				Status:              StatusInProgress,
				ScheduledFor:        nil,
			},
			expReq: &Request{
				ClientName:          baseClientName,
				ClientPhone:         baseClientPhone,
				Address:             baseAddress,
				ClientDescription:   baseClientDescription,
				EmployeeDescription: baseEmployeeDescription,
				CancelReason:        emptyCancelReason,
				Status:              StatusDone,
				ScheduledFor:        nil,
			},
			expErr: nil,
		},
		{
			name: "Попытка завершения новой заявки",
			req: &Request{
				ClientName:          baseClientName,
				ClientPhone:         baseClientPhone,
				Address:             baseAddress,
				ClientDescription:   baseClientDescription,
				EmployeeDescription: "",
				CancelReason:        emptyCancelReason,
				Status:              StatusNew,
				ScheduledFor:        nil,
			},
			expReq: &Request{
				ClientName:          baseClientName,
				ClientPhone:         baseClientPhone,
				Address:             baseAddress,
				ClientDescription:   baseClientDescription,
				EmployeeDescription: "",
				CancelReason:        emptyCancelReason,
				Status:              StatusNew,
				ScheduledFor:        nil,
			},
			expErr: ErrActionNotPermittedByStatus,
		},
		{
			name: "Попытка завершения отменённой заявки",
			req: &Request{
				ClientName:          baseClientName,
				ClientPhone:         baseClientPhone,
				Address:             baseAddress,
				ClientDescription:   baseClientDescription,
				EmployeeDescription: baseEmployeeDescription,
				CancelReason:        emptyCancelReason,
				Status:              StatusScheduled,
				ScheduledFor:        &threeDaysLater,
			},
			expReq: &Request{
				ClientName:          baseClientName,
				ClientPhone:         baseClientPhone,
				Address:             baseAddress,
				ClientDescription:   baseClientDescription,
				EmployeeDescription: baseEmployeeDescription,
				CancelReason:        filledCancelReason,
				Status:              StatusCanceled,
				ScheduledFor:        nil,
			},
			expErr: ErrActionNotPermittedByStatus,
		},
		{
			name: "Попытка завершения завершённой заявки",
			req: &Request{
				ClientName:          baseClientName,
				ClientPhone:         baseClientPhone,
				Address:             baseAddress,
				ClientDescription:   baseClientDescription,
				EmployeeDescription: baseEmployeeDescription,
				CancelReason:        emptyCancelReason,
				Status:              StatusDone,
				ScheduledFor:        nil,
			},
			expReq: &Request{
				ClientName:          baseClientName,
				ClientPhone:         baseClientPhone,
				Address:             baseAddress,
				ClientDescription:   baseClientDescription,
				EmployeeDescription: baseEmployeeDescription,
				CancelReason:        emptyCancelReason,
				Status:              StatusDone,
				ScheduledFor:        nil,
			},
			expErr: ErrActionNotPermittedByStatus,
		},
		{
			name: "Попытка завершения закрытой заявки",
			req: &Request{
				ClientName:          baseClientName,
				ClientPhone:         baseClientPhone,
				Address:             baseAddress,
				ClientDescription:   baseClientDescription,
				EmployeeDescription: baseEmployeeDescription,
				CancelReason:        emptyCancelReason,
				Status:              StatusPaid,
				ScheduledFor:        nil,
			},
			expReq: &Request{
				ClientName:          baseClientName,
				ClientPhone:         baseClientPhone,
				Address:             baseAddress,
				ClientDescription:   baseClientDescription,
				EmployeeDescription: baseEmployeeDescription,
				CancelReason:        emptyCancelReason,
				Status:              StatusPaid,
				ScheduledFor:        nil,
			},
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
			req: &Request{
				ClientName:          baseClientName,
				ClientPhone:         baseClientPhone,
				Address:             baseAddress,
				ClientDescription:   baseClientDescription,
				EmployeeDescription: baseEmployeeDescription,
				CancelReason:        emptyCancelReason,
				Status:              StatusDone,
				ScheduledFor:        nil,
			},
			expReq: &Request{
				ClientName:          baseClientName,
				ClientPhone:         baseClientPhone,
				Address:             baseAddress,
				ClientDescription:   baseClientDescription,
				EmployeeDescription: baseEmployeeDescription,
				CancelReason:        emptyCancelReason,
				Status:              StatusPaid,
				ScheduledFor:        nil,
			},
			expErr: nil,
		},
		{
			name: "Попытка закрытия назначенной заявки",
			req: &Request{
				ClientName:          baseClientName,
				ClientPhone:         baseClientPhone,
				Address:             baseAddress,
				ClientDescription:   baseClientDescription,
				EmployeeDescription: baseEmployeeDescription,
				CancelReason:        emptyCancelReason,
				Status:              StatusScheduled,
				ScheduledFor:        &threeDaysLater,
			},
			expReq: &Request{
				ClientName:          baseClientName,
				ClientPhone:         baseClientPhone,
				Address:             baseAddress,
				ClientDescription:   baseClientDescription,
				EmployeeDescription: baseEmployeeDescription,
				CancelReason:        emptyCancelReason,
				Status:              StatusScheduled,
				ScheduledFor:        &threeDaysLater,
			},
			expErr: ErrActionNotPermittedByStatus,
		},
		{
			name: "Попытка закрытия отменённой заявки",
			req: &Request{
				ClientName:          baseClientName,
				ClientPhone:         baseClientPhone,
				Address:             baseAddress,
				ClientDescription:   baseClientDescription,
				EmployeeDescription: baseEmployeeDescription,
				CancelReason:        emptyCancelReason,
				Status:              StatusScheduled,
				ScheduledFor:        &threeDaysLater,
			},
			expReq: &Request{
				ClientName:          baseClientName,
				ClientPhone:         baseClientPhone,
				Address:             baseAddress,
				ClientDescription:   baseClientDescription,
				EmployeeDescription: baseEmployeeDescription,
				CancelReason:        filledCancelReason,
				Status:              StatusCanceled,
				ScheduledFor:        nil,
			},
			expErr: ErrActionNotPermittedByStatus,
		},
		{
			name: "Попытка закрытия закрытой заявки",
			req: &Request{
				ClientName:          baseClientName,
				ClientPhone:         baseClientPhone,
				Address:             baseAddress,
				ClientDescription:   baseClientDescription,
				EmployeeDescription: baseEmployeeDescription,
				CancelReason:        emptyCancelReason,
				Status:              StatusPaid,
				ScheduledFor:        nil,
			},
			expReq: &Request{
				ClientName:          baseClientName,
				ClientPhone:         baseClientPhone,
				Address:             baseAddress,
				ClientDescription:   baseClientDescription,
				EmployeeDescription: baseEmployeeDescription,
				CancelReason:        emptyCancelReason,
				Status:              StatusPaid,
				ScheduledFor:        nil,
			},
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
			req: &Request{
				ClientName:          baseClientName,
				ClientPhone:         baseClientPhone,
				Address:             baseAddress,
				ClientDescription:   baseClientDescription,
				EmployeeDescription: "",
				CancelReason:        emptyCancelReason,
				Status:              StatusNew,
				ScheduledFor:        nil,
			},
			cancelReason: filledCancelReason,
			expReq: &Request{
				ClientName:          baseClientName,
				ClientPhone:         baseClientPhone,
				Address:             baseAddress,
				ClientDescription:   baseClientDescription,
				EmployeeDescription: baseEmployeeDescription,
				CancelReason:        filledCancelReason,
				Status:              StatusCanceled,
				ScheduledFor:        nil,
			},
			expErr: nil,
		},
		{
			name: "Успешная отмена запланированной заявки",
			req: &Request{
				ClientName:          baseClientName,
				ClientPhone:         baseClientPhone,
				Address:             baseAddress,
				ClientDescription:   baseClientDescription,
				EmployeeDescription: baseEmployeeDescription,
				CancelReason:        emptyCancelReason,
				Status:              StatusScheduled,
				ScheduledFor:        &threeDaysLater,
			},
			cancelReason: filledCancelReason,
			expReq: &Request{
				ClientName:          baseClientName,
				ClientPhone:         baseClientPhone,
				Address:             baseAddress,
				ClientDescription:   baseClientDescription,
				EmployeeDescription: baseEmployeeDescription,
				CancelReason:        filledCancelReason,
				Status:              StatusCanceled,
				ScheduledFor:        nil,
			},
			expErr: nil,
		},
		{
			name: "Попытка отмены отменённой заявки",
			req: &Request{
				ClientName:          baseClientName,
				ClientPhone:         baseClientPhone,
				Address:             baseAddress,
				ClientDescription:   baseClientDescription,
				EmployeeDescription: baseEmployeeDescription,
				CancelReason:        filledCancelReason,
				Status:              StatusCanceled,
				ScheduledFor:        nil,
			},
			cancelReason: "Другая причина отмены",
			expReq: &Request{
				ClientName:          baseClientName,
				ClientPhone:         baseClientPhone,
				Address:             baseAddress,
				ClientDescription:   baseClientDescription,
				EmployeeDescription: baseEmployeeDescription,
				CancelReason:        filledCancelReason,
				Status:              StatusCanceled,
				ScheduledFor:        nil,
			},
			expErr: ErrActionNotPermittedByStatus,
		},
		{
			name: "Попытка отмены оплаченной заявки",
			req: &Request{
				ClientName:          baseClientName,
				ClientPhone:         baseClientPhone,
				Address:             baseAddress,
				ClientDescription:   baseClientDescription,
				EmployeeDescription: baseEmployeeDescription,
				CancelReason:        emptyCancelReason,
				Status:              StatusPaid,
				ScheduledFor:        nil,
			},
			cancelReason: filledCancelReason,
			expReq: &Request{
				ClientName:          baseClientName,
				ClientPhone:         baseClientPhone,
				Address:             baseAddress,
				ClientDescription:   baseClientDescription,
				EmployeeDescription: baseEmployeeDescription,
				CancelReason:        emptyCancelReason,
				Status:              StatusPaid,
				ScheduledFor:        nil,
			},
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
			req: &Request{
				ClientName:          baseClientName,
				ClientPhone:         baseClientPhone,
				Address:             baseAddress,
				ClientDescription:   baseClientDescription,
				EmployeeDescription: baseEmployeeDescription,
				CancelReason:        emptyCancelReason,
				Status:              StatusScheduled,
				ScheduledFor:        &threeDaysLater,
			},
			patchedFields: &RequestPatcher{
				ClientName:          &patchedClientName,
				ClientPhone:         &patchedClientPhone,
				Address:             &patchedAddress,
				ClientDescription:   &patchedCliendDescription,
				EmployeeDescription: &patchedEmployeeDescription,
			},
			expReq: &Request{
				ClientName:          patchedClientName,
				ClientPhone:         patchedClientPhone,
				Address:             patchedAddress,
				ClientDescription:   patchedCliendDescription,
				EmployeeDescription: patchedEmployeeDescription,
				CancelReason:        emptyCancelReason,
				Status:              StatusScheduled,
				ScheduledFor:        nil,
			},
			expErr: nil,
		},
		{
			name: "Попытка модификации отменённой заявки",
			req: &Request{
				ClientName:          baseClientName,
				ClientPhone:         baseClientPhone,
				Address:             baseAddress,
				ClientDescription:   baseClientDescription,
				EmployeeDescription: baseEmployeeDescription,
				CancelReason:        filledCancelReason,
				Status:              StatusCanceled,
				ScheduledFor:        nil,
			},
			patchedFields: &RequestPatcher{
				ClientName:          &patchedClientName,
				ClientPhone:         &patchedClientPhone,
				Address:             &patchedAddress,
				ClientDescription:   &patchedCliendDescription,
				EmployeeDescription: &patchedEmployeeDescription,
			},
			expReq: &Request{
				ClientName:          baseClientName,
				ClientPhone:         baseClientPhone,
				Address:             baseAddress,
				ClientDescription:   baseClientDescription,
				EmployeeDescription: baseEmployeeDescription,
				CancelReason:        filledCancelReason,
				Status:              StatusCanceled,
				ScheduledFor:        nil,
			},
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
