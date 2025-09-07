package requests_test

import (
	"testing"
	"time"

	"github.com/seagumineko/spkuznetsov/internal/core/logic/auth"
	"github.com/seagumineko/spkuznetsov/internal/core/logic/requests"
	core_errors "github.com/seagumineko/spkuznetsov/internal/errors"
	"github.com/seagumineko/spkuznetsov/internal/testutils"
)

var (
	tomorrow       = testutils.GetNDaysLater(1)
	threeDaysLater = testutils.GetNDaysLater(3)
)

func TestNewRequest(t *testing.T) {
	basePrimaryRequest := &requests.PrimaryRequest{
		ClientName:        testutils.ClientName,
		ClientPhone:       testutils.ClientPhone,
		Address:           testutils.Address,
		ClientDescription: testutils.ClientDescription,
	}

	cases := []struct {
		name   string
		pReq   *requests.PrimaryRequest
		expReq *requests.Request
		expErr error
	}{
		{
			name: "Успешное создание новой заявки",
			pReq: basePrimaryRequest,
			expReq: testutils.NewTestRequest(
				testutils.WithClientName(basePrimaryRequest.ClientName),
				testutils.WithClientPhone(basePrimaryRequest.ClientPhone),
				testutils.WithAddress(basePrimaryRequest.Address),
				testutils.WithClientDescription(basePrimaryRequest.ClientDescription),
				testutils.WithEmployee(nil),
				testutils.WithStatus(requests.StatusNew),
				testutils.WithScheduledFor(nil),
			),
			expErr: nil,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req, err := c.pReq.CreateNewRequest()
			testutils.AssertError(t, c.expErr, err)
			testutils.ValidateRequest(t, c.expReq, req)
		})
	}
}

func TestPreschedule(t *testing.T) {
	cases := []struct {
		name   string
		req    *requests.Request
		date   *time.Time
		expReq *requests.Request
		expErr error
	}{
		{
			name: "Успешная попытка назначить предварительную дату новой заявки",
			req: testutils.NewTestRequest(
				testutils.WithStatus(requests.StatusNew),
				testutils.WithScheduledFor(nil),
			),
			date: &tomorrow,
			expReq: testutils.NewTestRequest(
				testutils.WithStatus(requests.StatusPrescheduled),
				testutils.WithScheduledFor(&tomorrow),
			),
			expErr: nil,
		},
		{
			name: "Успешная попытка переназначить предварительную дату новой заявки",
			req: testutils.NewTestRequest(
				testutils.WithStatus(requests.StatusPrescheduled),
				testutils.WithScheduledFor(&tomorrow),
			),
			date: &threeDaysLater,
			expReq: testutils.NewTestRequest(
				testutils.WithStatus(requests.StatusPrescheduled),
				testutils.WithScheduledFor(&threeDaysLater),
			),
			expErr: nil,
		},
		{
			name: "Успешная попытка назначить предварительную дату заявки после частичных работ",
			req: testutils.NewTestRequest(
				testutils.WithStatus(requests.StatusInProgress),
				testutils.WithScheduledFor(nil),
			),
			date: &tomorrow,
			expReq: testutils.NewTestRequest(
				testutils.WithStatus(requests.StatusPrescheduled),
				testutils.WithScheduledFor(&tomorrow),
			),
			expErr: nil,
		},
		{
			name: "Попытка назначить предварительную дату для отменённой заявки",
			req: testutils.NewTestRequest(
				testutils.WithStatus(requests.StatusCanceled),
				testutils.WithScheduledFor(nil),
			),
			date: &tomorrow,
			expReq: testutils.NewTestRequest(
				testutils.WithStatus(requests.StatusCanceled),
				testutils.WithScheduledFor(nil),
			),
			expErr: core_errors.ErrRequestActionNotPermittedByStatus,
		},
		{
			name: "Попытка назначить предварительную дату для выполненной заявки",
			req: testutils.NewTestRequest(
				testutils.WithStatus(requests.StatusDone),
				testutils.WithScheduledFor(nil),
			),
			date: &tomorrow,
			expReq: testutils.NewTestRequest(
				testutils.WithStatus(requests.StatusDone),
				testutils.WithScheduledFor(nil),
			),
			expErr: core_errors.ErrRequestActionNotPermittedByStatus,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.req.Preschedule(c.date)
			testutils.AssertError(t, c.expErr, err)
			testutils.ValidateRequest(t, c.expReq, c.req)
		})
	}
}

func TestAssign(t *testing.T) {
	employee := &auth.Employee{
		Name: "Николай Николаев",
	}
	cases := []struct {
		name   string
		req    *requests.Request
		emp    *auth.Employee
		expReq *requests.Request
		expErr error
	}{
		{
			name: "Успешная попытка назначить сотрудника на новую заявку",
			req: testutils.NewTestRequest(
				testutils.WithEmployee(nil),
				testutils.WithStatus(requests.StatusPrescheduled),
			),
			emp: employee,
			expReq: testutils.NewTestRequest(
				testutils.WithEmployee(employee),
				testutils.WithStatus(requests.StatusAssigned),
			),
			expErr: nil,
		},
		{
			name: "Попытка назначить сотрудника на отменённую заявку",
			req: testutils.NewTestRequest(
				testutils.WithEmployee(nil),
				testutils.WithStatus(requests.StatusCanceled),
			),
			emp: employee,
			expReq: testutils.NewTestRequest(
				testutils.WithEmployee(nil),
				testutils.WithStatus(requests.StatusCanceled),
			),
			expErr: core_errors.ErrRequestActionNotPermittedByStatus,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.req.Assign(c.emp)
			testutils.AssertError(t, c.expErr, err)
			testutils.ValidateRequest(t, c.expReq, c.req)
		})
	}
}

func TestSchedule(t *testing.T) {
	cases := []struct {
		name     string
		req      *requests.Request
		schedule *time.Time
		expReq   *requests.Request
		expErr   error
	}{
		{
			name: "Успешное планирование даты работ",
			req: testutils.NewTestRequest(
				testutils.WithStatus(requests.StatusAssigned),
				testutils.WithScheduledFor(&threeDaysLater),
			),
			schedule: &tomorrow,
			expReq: testutils.NewTestRequest(
				testutils.WithStatus(requests.StatusScheduled),
				testutils.WithScheduledFor(&tomorrow),
			),
			expErr: nil,
		},
		{
			name: "Успешное планирование даты новых работ",
			req: testutils.NewTestRequest(
				testutils.WithStatus(requests.StatusInProgress),
				testutils.WithScheduledFor(nil),
			),
			schedule: &tomorrow,
			expReq: testutils.NewTestRequest(
				testutils.WithStatus(requests.StatusScheduled),
				testutils.WithScheduledFor(&tomorrow),
			),
			expErr: nil,
		},
		{
			name: "Попытка запланировать выполненные даты работы",
			req: testutils.NewTestRequest(
				testutils.WithStatus(requests.StatusDone),
			),
			schedule: &tomorrow,
			expReq: testutils.NewTestRequest(
				testutils.WithStatus(requests.StatusDone),
			),
			expErr: core_errors.ErrRequestActionNotPermittedByStatus,
		},
		{
			name: "Попытка запланировать отменённые работы",
			req: testutils.NewTestRequest(
				testutils.WithStatus(requests.StatusCanceled),
			),
			schedule: &threeDaysLater,
			expReq: testutils.NewTestRequest(
				testutils.WithStatus(requests.StatusCanceled),
			),
			expErr: core_errors.ErrRequestActionNotPermittedByStatus,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.req.Schedule(c.schedule)
			testutils.AssertError(t, c.expErr, err)
			testutils.ValidateRequest(t, c.expReq, c.req)
		})
	}
}

func TestConfirmSchedule(t *testing.T) {
	tomorrow := testutils.GetNDaysLater(1)
	cases := []struct {
		name   string
		req    *requests.Request
		expReq *requests.Request
		expErr error
	}{
		{
			name: "Успешное подтверждение даты работ",
			req: testutils.NewTestRequest(
				testutils.WithStatus(requests.StatusAssigned),
				testutils.WithScheduledFor(&tomorrow),
			),
			expReq: testutils.NewTestRequest(
				testutils.WithStatus(requests.StatusScheduled),
				testutils.WithScheduledFor(&tomorrow),
			),
			expErr: nil,
		},
		{
			name: "Попытка подтвердить работы без предварительной даты",
			req: testutils.NewTestRequest(
				testutils.WithStatus(requests.StatusAssigned),
				testutils.WithScheduledFor(nil),
			),
			expReq: testutils.NewTestRequest(
				testutils.WithStatus(requests.StatusAssigned),
				testutils.WithScheduledFor(nil),
			),
			expErr: core_errors.ErrInvalidDate,
		},
		{
			name: "Попытка подтвердить отменённые работы",
			req: testutils.NewTestRequest(
				testutils.WithStatus(requests.StatusCanceled),
			),
			expReq: testutils.NewTestRequest(
				testutils.WithStatus(requests.StatusCanceled),
			),
			expErr: core_errors.ErrRequestActionNotPermittedByStatus,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.req.ConfirmSchedule()
			testutils.AssertError(t, c.expErr, err)
			testutils.ValidateRequest(t, c.expReq, c.req)
		})
	}
}

func TestProgress(t *testing.T) {
	cases := []struct {
		name    string
		req     *requests.Request
		empDesc string
		expReq  *requests.Request
		expErr  error
	}{
		{
			name: "Успешный прогресс заявки",
			req: testutils.NewTestRequest(
				testutils.WithStatus(requests.StatusScheduled),
				testutils.WithScheduledFor(&tomorrow),
			),
			empDesc: testutils.EmployeeDescription,
			expReq: testutils.NewTestRequest(
				testutils.WithStatus(requests.StatusInProgress),
				testutils.WithScheduledFor(nil),
				testutils.WithEmployeeDescription(testutils.EmployeeDescription),
			),
			expErr: nil,
		},
		{
			name: "Попытка прогресса в выполненной заявке",
			req: testutils.NewTestRequest(
				testutils.WithStatus(requests.StatusDone),
			),
			empDesc: "Новое описание",
			expReq: testutils.NewTestRequest(
				testutils.WithStatus(requests.StatusDone),
			),
			expErr: core_errors.ErrRequestActionNotPermittedByStatus,
		},
		{
			name: "Попытка прогресса в отменённой заявке",
			req: testutils.NewTestRequest(
				testutils.WithStatus(requests.StatusCanceled),
			),
			empDesc: "Новое описание",
			expReq: testutils.NewTestRequest(
				testutils.WithStatus(requests.StatusCanceled),
			),
			expErr: core_errors.ErrRequestActionNotPermittedByStatus,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.req.Progress(c.empDesc)
			testutils.AssertError(t, c.expErr, err)
			testutils.ValidateRequest(t, c.expReq, c.req)
		})
	}
}

func TestComplete(t *testing.T) {
	cases := []struct {
		name   string
		req    *requests.Request
		expReq *requests.Request
		expErr error
	}{
		{
			name: "Успешное завершение заявки",
			req: testutils.NewTestRequest(
				testutils.WithStatus(requests.StatusInProgress),
			),
			expReq: testutils.NewTestRequest(
				testutils.WithStatus(requests.StatusDone),
			),
			expErr: nil,
		},
		{
			name: "Попытка завершения новой заявки",
			req: testutils.NewTestRequest(
				testutils.WithStatus(requests.StatusNew),
			),
			expReq: testutils.NewTestRequest(
				testutils.WithStatus(requests.StatusNew),
			),
			expErr: core_errors.ErrRequestActionNotPermittedByStatus,
		},
		{
			name: "Попытка завершения отменённой заявки",
			req: testutils.NewTestRequest(
				testutils.WithStatus(requests.StatusCanceled),
			),
			expReq: testutils.NewTestRequest(
				testutils.WithStatus(requests.StatusCanceled),
			),
			expErr: core_errors.ErrRequestActionNotPermittedByStatus,
		},
		{
			name: "Попытка завершения завершённой заявки",
			req: testutils.NewTestRequest(
				testutils.WithStatus(requests.StatusDone),
			),
			expReq: testutils.NewTestRequest(
				testutils.WithStatus(requests.StatusDone),
			),
			expErr: core_errors.ErrRequestActionNotPermittedByStatus,
		},
		{
			name: "Попытка завершения закрытой заявки",
			req: testutils.NewTestRequest(
				testutils.WithStatus(requests.StatusPaid),
			),
			expReq: testutils.NewTestRequest(
				testutils.WithStatus(requests.StatusPaid),
			),
			expErr: core_errors.ErrRequestActionNotPermittedByStatus,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.req.Complete()
			testutils.AssertError(t, c.expErr, err)
			testutils.ValidateRequest(t, c.expReq, c.req)
		})
	}
}

func TestClose(t *testing.T) {
	cases := []struct {
		name   string
		req    *requests.Request
		expReq *requests.Request
		expErr error
	}{
		{
			name: "Успешное закрытие выполненной заявки",
			req: testutils.NewTestRequest(
				testutils.WithStatus(requests.StatusDone),
			),
			expReq: testutils.NewTestRequest(
				testutils.WithStatus(requests.StatusPaid),
			),
			expErr: nil,
		},
		{
			name: "Попытка закрытия назначенной заявки",
			req: testutils.NewTestRequest(
				testutils.WithStatus(requests.StatusScheduled),
			),
			expReq: testutils.NewTestRequest(
				testutils.WithStatus(requests.StatusScheduled),
			),
			expErr: core_errors.ErrRequestActionNotPermittedByStatus,
		},
		{
			name: "Попытка закрытия отменённой заявки",
			req: testutils.NewTestRequest(
				testutils.WithStatus(requests.StatusCanceled),
			),
			expReq: testutils.NewTestRequest(
				testutils.WithStatus(requests.StatusCanceled),
			),
			expErr: core_errors.ErrRequestActionNotPermittedByStatus,
		},
		{
			name: "Попытка закрытия закрытой заявки",
			req: testutils.NewTestRequest(
				testutils.WithStatus(requests.StatusPaid),
			),
			expReq: testutils.NewTestRequest(
				testutils.WithStatus(requests.StatusPaid),
			),
			expErr: core_errors.ErrRequestActionNotPermittedByStatus,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.req.Close()
			testutils.AssertError(t, c.expErr, err)
			testutils.ValidateRequest(t, c.expReq, c.req)
		})
	}
}

func TestCancel(t *testing.T) {
	cases := []struct {
		name         string
		req          *requests.Request
		cancelReason string
		expReq       *requests.Request
		expErr       error
	}{
		{
			name: "Успешная отмена новой заявки",
			req: testutils.NewTestRequest(
				testutils.WithStatus(requests.StatusNew),
			),
			cancelReason: testutils.FilledCancelReason,
			expReq: testutils.NewTestRequest(
				testutils.WithStatus(requests.StatusCanceled),
				testutils.WithCancelReason(testutils.FilledCancelReason),
				testutils.WithScheduledFor(nil),
			),
			expErr: nil,
		},
		{
			name: "Успешная отмена запланированной заявки",
			req: testutils.NewTestRequest(
				testutils.WithStatus(requests.StatusScheduled),
			),
			cancelReason: testutils.FilledCancelReason,
			expReq: testutils.NewTestRequest(
				testutils.WithStatus(requests.StatusCanceled),
				testutils.WithCancelReason(testutils.FilledCancelReason),
				testutils.WithScheduledFor(nil),
			),
			expErr: nil,
		},
		{
			name: "Попытка отмены отменённой заявки",
			req: testutils.NewTestRequest(
				testutils.WithStatus(requests.StatusCanceled),
				testutils.WithCancelReason(testutils.FilledCancelReason),
			),
			cancelReason: "Другая причина отмены",
			expReq: testutils.NewTestRequest(
				testutils.WithStatus(requests.StatusCanceled),
				testutils.WithCancelReason(testutils.FilledCancelReason),
			),
			expErr: core_errors.ErrRequestActionNotPermittedByStatus,
		},
		{
			name: "Попытка отмены оплаченной заявки",
			req: testutils.NewTestRequest(
				testutils.WithStatus(requests.StatusPaid),
			),
			cancelReason: testutils.FilledCancelReason,
			expReq: testutils.NewTestRequest(
				testutils.WithStatus(requests.StatusPaid),
			),
			expErr: core_errors.ErrRequestActionNotPermittedByStatus,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.req.Cancel(c.cancelReason)
			testutils.AssertError(t, c.expErr, err)
			testutils.ValidateRequest(t, c.expReq, c.req)
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
		req           *requests.Request
		patchedFields *requests.RequestPatcher
		expReq        *requests.Request
		expErr        error
	}{
		{
			name: "Успешная модификация полей заявки",
			req:  testutils.NewTestRequest(),
			patchedFields: &requests.RequestPatcher{
				ClientName:          &patchedClientName,
				ClientPhone:         &patchedClientPhone,
				Address:             &patchedAddress,
				ClientDescription:   &patchedCliendDescription,
				EmployeeDescription: &patchedEmployeeDescription,
			},
			expReq: testutils.NewTestRequest(
				testutils.WithClientName(patchedClientName),
				testutils.WithClientPhone(patchedClientPhone),
				testutils.WithAddress(patchedAddress),
				testutils.WithClientDescription(patchedCliendDescription),
				testutils.WithEmployeeDescription(patchedEmployeeDescription),
			),
			expErr: nil,
		},
		{
			name: "Попытка модификации отменённой заявки",
			req: testutils.NewTestRequest(
				testutils.WithStatus(requests.StatusCanceled),
			),
			patchedFields: &requests.RequestPatcher{
				ClientName:          &patchedClientName,
				ClientPhone:         &patchedClientPhone,
				Address:             &patchedAddress,
				ClientDescription:   &patchedCliendDescription,
				EmployeeDescription: &patchedEmployeeDescription,
			},
			expReq: testutils.NewTestRequest(
				testutils.WithStatus(requests.StatusCanceled),
			),
			expErr: core_errors.ErrRequestActionNotPermittedByStatus,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.req.Patch(c.patchedFields)
			testutils.AssertError(t, c.expErr, err)
			testutils.ValidateRequest(t, c.expReq, c.req)
		})
	}
}
