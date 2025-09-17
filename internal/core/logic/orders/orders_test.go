package orders_test

import (
	"testing"
	"time"

	"github.com/Owouwun/spkuznetsov/internal/core/logic/auth"
	"github.com/Owouwun/spkuznetsov/internal/core/logic/orders"
	deterrs "github.com/Owouwun/spkuznetsov/internal/errors"
	"github.com/Owouwun/spkuznetsov/internal/testutils"
)

var (
	tomorrow       = testutils.GetNDaysLater(1)
	threeDaysLater = testutils.GetNDaysLater(3)
)

func TestNewOrder(t *testing.T) {
	basePrimaryOrder := &orders.PrimaryOrder{
		ClientName:        testutils.ClientName,
		ClientPhone:       testutils.ClientPhone,
		Address:           testutils.Address,
		ClientDescription: testutils.ClientDescription,
	}

	cases := []struct {
		name   string
		pReq   *orders.PrimaryOrder
		expReq *orders.Order
		expErr error
	}{
		{
			name: "Успешное создание новой заявки",
			pReq: basePrimaryOrder,
			expReq: testutils.NewTestOrder(
				testutils.WithClientName(basePrimaryOrder.ClientName),
				testutils.WithClientPhone(basePrimaryOrder.ClientPhone),
				testutils.WithAddress(basePrimaryOrder.Address),
				testutils.WithClientDescription(basePrimaryOrder.ClientDescription),
				testutils.WithEmployee(nil),
				testutils.WithStatus(orders.StatusNew),
				testutils.WithScheduledFor(nil),
			),
			expErr: nil,
		},
		{
			name: "Попытка создать заявку с некорректным номером телефона клиента",
			pReq: &orders.PrimaryOrder{
				ClientName:        testutils.ClientName,
				ClientPhone:       "+7111222334",
				Address:           testutils.Address,
				ClientDescription: testutils.ClientDescription,
			},
			expReq: nil,
			expErr: deterrs.NewDetErr(
				deterrs.InvalidValue,
			),
		},
		{
			name: "Попытка создать заявку без имени клиента",
			pReq: &orders.PrimaryOrder{
				ClientName:        "",
				ClientPhone:       testutils.ClientPhone,
				Address:           testutils.Address,
				ClientDescription: testutils.ClientDescription,
			},
			expReq: nil,
			expErr: deterrs.NewDetErr(
				deterrs.EmptyField,
			),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req, err := c.pReq.CreateNewOrder()
			testutils.AssertError(t, c.expErr, err)
			testutils.ValidateOrder(t, c.expReq, req)
		})
	}
}

func TestPreschedule(t *testing.T) {
	cases := []struct {
		name   string
		req    *orders.Order
		date   *time.Time
		expReq *orders.Order
		expErr error
	}{
		{
			name: "Успешная попытка назначить предварительную дату новой заявки",
			req: testutils.NewTestOrder(
				testutils.WithStatus(orders.StatusNew),
				testutils.WithScheduledFor(nil),
			),
			date: &tomorrow,
			expReq: testutils.NewTestOrder(
				testutils.WithStatus(orders.StatusPrescheduled),
				testutils.WithScheduledFor(&tomorrow),
			),
			expErr: nil,
		},
		{
			name: "Успешная попытка переназначить предварительную дату новой заявки",
			req: testutils.NewTestOrder(
				testutils.WithStatus(orders.StatusPrescheduled),
				testutils.WithScheduledFor(&tomorrow),
			),
			date: &threeDaysLater,
			expReq: testutils.NewTestOrder(
				testutils.WithStatus(orders.StatusPrescheduled),
				testutils.WithScheduledFor(&threeDaysLater),
			),
			expErr: nil,
		},
		{
			name: "Успешная попытка назначить предварительную дату заявки после частичных работ",
			req: testutils.NewTestOrder(
				testutils.WithStatus(orders.StatusInProgress),
				testutils.WithScheduledFor(nil),
			),
			date: &tomorrow,
			expReq: testutils.NewTestOrder(
				testutils.WithStatus(orders.StatusPrescheduled),
				testutils.WithScheduledFor(&tomorrow),
			),
			expErr: nil,
		},
		{
			name: "Попытка назначить предварительную дату для отменённой заявки",
			req: testutils.NewTestOrder(
				testutils.WithStatus(orders.StatusCanceled),
				testutils.WithScheduledFor(nil),
			),
			date: &tomorrow,
			expReq: testutils.NewTestOrder(
				testutils.WithStatus(orders.StatusCanceled),
				testutils.WithScheduledFor(nil),
			),
			expErr: deterrs.NewDetErr(
				deterrs.OrderActionNotPermittedByStatus,
			),
		},
		{
			name: "Попытка назначить предварительную дату для выполненной заявки",
			req: testutils.NewTestOrder(
				testutils.WithStatus(orders.StatusDone),
				testutils.WithScheduledFor(nil),
			),
			date: &tomorrow,
			expReq: testutils.NewTestOrder(
				testutils.WithStatus(orders.StatusDone),
				testutils.WithScheduledFor(nil),
			),
			expErr: deterrs.NewDetErr(
				deterrs.OrderActionNotPermittedByStatus,
			),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.req.Preschedule(c.date)
			testutils.AssertError(t, c.expErr, err)
			testutils.ValidateOrder(t, c.expReq, c.req)
		})
	}
}

func TestAssign(t *testing.T) {
	employee := &auth.Employee{
		Name: "Николай Николаев",
	}
	cases := []struct {
		name   string
		req    *orders.Order
		emp    *auth.Employee
		expReq *orders.Order
		expErr error
	}{
		{
			name: "Успешная попытка назначить сотрудника на новую заявку",
			req: testutils.NewTestOrder(
				testutils.WithEmployee(nil),
				testutils.WithStatus(orders.StatusPrescheduled),
			),
			emp: employee,
			expReq: testutils.NewTestOrder(
				testutils.WithEmployee(employee),
				testutils.WithStatus(orders.StatusAssigned),
			),
			expErr: nil,
		},
		{
			name: "Попытка назначить сотрудника на отменённую заявку",
			req: testutils.NewTestOrder(
				testutils.WithEmployee(nil),
				testutils.WithStatus(orders.StatusCanceled),
			),
			emp: employee,
			expReq: testutils.NewTestOrder(
				testutils.WithEmployee(nil),
				testutils.WithStatus(orders.StatusCanceled),
			),
			expErr: deterrs.NewDetErr(
				deterrs.OrderActionNotPermittedByStatus,
			),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.req.Assign(c.emp)
			testutils.AssertError(t, c.expErr, err)
			testutils.ValidateOrder(t, c.expReq, c.req)
		})
	}
}

func TestSchedule(t *testing.T) {
	cases := []struct {
		name     string
		req      *orders.Order
		schedule *time.Time
		expReq   *orders.Order
		expErr   error
	}{
		{
			name: "Успешное планирование даты работ",
			req: testutils.NewTestOrder(
				testutils.WithStatus(orders.StatusAssigned),
				testutils.WithScheduledFor(&threeDaysLater),
			),
			schedule: &tomorrow,
			expReq: testutils.NewTestOrder(
				testutils.WithStatus(orders.StatusScheduled),
				testutils.WithScheduledFor(&tomorrow),
			),
			expErr: nil,
		},
		{
			name: "Успешное планирование даты новых работ",
			req: testutils.NewTestOrder(
				testutils.WithStatus(orders.StatusInProgress),
				testutils.WithScheduledFor(nil),
			),
			schedule: &tomorrow,
			expReq: testutils.NewTestOrder(
				testutils.WithStatus(orders.StatusScheduled),
				testutils.WithScheduledFor(&tomorrow),
			),
			expErr: nil,
		},
		{
			name: "Попытка запланировать выполненные даты работы",
			req: testutils.NewTestOrder(
				testutils.WithStatus(orders.StatusDone),
			),
			schedule: &tomorrow,
			expReq: testutils.NewTestOrder(
				testutils.WithStatus(orders.StatusDone),
			),
			expErr: deterrs.NewDetErr(
				deterrs.OrderActionNotPermittedByStatus,
			),
		},
		{
			name: "Попытка запланировать отменённые работы",
			req: testutils.NewTestOrder(
				testutils.WithStatus(orders.StatusCanceled),
			),
			schedule: &threeDaysLater,
			expReq: testutils.NewTestOrder(
				testutils.WithStatus(orders.StatusCanceled),
			),
			expErr: deterrs.NewDetErr(
				deterrs.OrderActionNotPermittedByStatus,
			),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.req.Schedule(c.schedule)
			testutils.AssertError(t, c.expErr, err)
			testutils.ValidateOrder(t, c.expReq, c.req)
		})
	}
}

func TestConfirmSchedule(t *testing.T) {
	tomorrow := testutils.GetNDaysLater(1)
	cases := []struct {
		name   string
		req    *orders.Order
		expReq *orders.Order
		expErr error
	}{
		{
			name: "Успешное подтверждение даты работ",
			req: testutils.NewTestOrder(
				testutils.WithStatus(orders.StatusAssigned),
				testutils.WithScheduledFor(&tomorrow),
			),
			expReq: testutils.NewTestOrder(
				testutils.WithStatus(orders.StatusScheduled),
				testutils.WithScheduledFor(&tomorrow),
			),
			expErr: nil,
		},
		{
			name: "Попытка подтвердить работы без предварительной даты",
			req: testutils.NewTestOrder(
				testutils.WithStatus(orders.StatusAssigned),
				testutils.WithScheduledFor(nil),
			),
			expReq: testutils.NewTestOrder(
				testutils.WithStatus(orders.StatusAssigned),
				testutils.WithScheduledFor(nil),
			),
			expErr: deterrs.NewDetErr(
				deterrs.EmptyField,
			),
		},
		{
			name: "Попытка подтвердить отменённые работы",
			req: testutils.NewTestOrder(
				testutils.WithStatus(orders.StatusCanceled),
			),
			expReq: testutils.NewTestOrder(
				testutils.WithStatus(orders.StatusCanceled),
			),
			expErr: deterrs.NewDetErr(
				deterrs.OrderActionNotPermittedByStatus,
			),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.req.ConfirmSchedule()
			testutils.AssertError(t, c.expErr, err)
			testutils.ValidateOrder(t, c.expReq, c.req)
		})
	}
}

func TestProgress(t *testing.T) {
	cases := []struct {
		name    string
		req     *orders.Order
		empDesc string
		expReq  *orders.Order
		expErr  error
	}{
		{
			name: "Успешный прогресс заявки",
			req: testutils.NewTestOrder(
				testutils.WithStatus(orders.StatusScheduled),
				testutils.WithScheduledFor(&tomorrow),
			),
			empDesc: testutils.EmployeeDescription,
			expReq: testutils.NewTestOrder(
				testutils.WithStatus(orders.StatusInProgress),
				testutils.WithScheduledFor(nil),
				testutils.WithEmployeeDescription(testutils.EmployeeDescription),
			),
			expErr: nil,
		},
		{
			name: "Попытка прогресса в выполненной заявке",
			req: testutils.NewTestOrder(
				testutils.WithStatus(orders.StatusDone),
			),
			empDesc: "Новое описание",
			expReq: testutils.NewTestOrder(
				testutils.WithStatus(orders.StatusDone),
			),
			expErr: deterrs.NewDetErr(
				deterrs.OrderActionNotPermittedByStatus,
			),
		},
		{
			name: "Попытка прогресса в отменённой заявке",
			req: testutils.NewTestOrder(
				testutils.WithStatus(orders.StatusCanceled),
			),
			empDesc: "Новое описание",
			expReq: testutils.NewTestOrder(
				testutils.WithStatus(orders.StatusCanceled),
			),
			expErr: deterrs.NewDetErr(
				deterrs.OrderActionNotPermittedByStatus,
			),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.req.Progress(c.empDesc)
			testutils.AssertError(t, c.expErr, err)
			testutils.ValidateOrder(t, c.expReq, c.req)
		})
	}
}

func TestComplete(t *testing.T) {
	cases := []struct {
		name   string
		req    *orders.Order
		expReq *orders.Order
		expErr error
	}{
		{
			name: "Успешное завершение заявки",
			req: testutils.NewTestOrder(
				testutils.WithStatus(orders.StatusInProgress),
			),
			expReq: testutils.NewTestOrder(
				testutils.WithStatus(orders.StatusDone),
			),
			expErr: nil,
		},
		{
			name: "Попытка завершения новой заявки",
			req: testutils.NewTestOrder(
				testutils.WithStatus(orders.StatusNew),
			),
			expReq: testutils.NewTestOrder(
				testutils.WithStatus(orders.StatusNew),
			),
			expErr: deterrs.NewDetErr(
				deterrs.OrderActionNotPermittedByStatus,
			),
		},
		{
			name: "Попытка завершения отменённой заявки",
			req: testutils.NewTestOrder(
				testutils.WithStatus(orders.StatusCanceled),
			),
			expReq: testutils.NewTestOrder(
				testutils.WithStatus(orders.StatusCanceled),
			),
			expErr: deterrs.NewDetErr(
				deterrs.OrderActionNotPermittedByStatus,
			),
		},
		{
			name: "Попытка завершения завершённой заявки",
			req: testutils.NewTestOrder(
				testutils.WithStatus(orders.StatusDone),
			),
			expReq: testutils.NewTestOrder(
				testutils.WithStatus(orders.StatusDone),
			),
			expErr: deterrs.NewDetErr(
				deterrs.OrderActionNotPermittedByStatus,
			),
		},
		{
			name: "Попытка завершения закрытой заявки",
			req: testutils.NewTestOrder(
				testutils.WithStatus(orders.StatusPaid),
			),
			expReq: testutils.NewTestOrder(
				testutils.WithStatus(orders.StatusPaid),
			),
			expErr: deterrs.NewDetErr(
				deterrs.OrderActionNotPermittedByStatus,
			),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.req.Complete()
			testutils.AssertError(t, c.expErr, err)
			testutils.ValidateOrder(t, c.expReq, c.req)
		})
	}
}

func TestClose(t *testing.T) {
	cases := []struct {
		name   string
		req    *orders.Order
		expReq *orders.Order
		expErr error
	}{
		{
			name: "Успешное закрытие выполненной заявки",
			req: testutils.NewTestOrder(
				testutils.WithStatus(orders.StatusDone),
			),
			expReq: testutils.NewTestOrder(
				testutils.WithStatus(orders.StatusPaid),
			),
			expErr: nil,
		},
		{
			name: "Попытка закрытия назначенной заявки",
			req: testutils.NewTestOrder(
				testutils.WithStatus(orders.StatusScheduled),
			),
			expReq: testutils.NewTestOrder(
				testutils.WithStatus(orders.StatusScheduled),
			),
			expErr: deterrs.NewDetErr(
				deterrs.OrderActionNotPermittedByStatus,
			),
		},
		{
			name: "Попытка закрытия отменённой заявки",
			req: testutils.NewTestOrder(
				testutils.WithStatus(orders.StatusCanceled),
			),
			expReq: testutils.NewTestOrder(
				testutils.WithStatus(orders.StatusCanceled),
			),
			expErr: deterrs.NewDetErr(
				deterrs.OrderActionNotPermittedByStatus,
			),
		},
		{
			name: "Попытка закрытия закрытой заявки",
			req: testutils.NewTestOrder(
				testutils.WithStatus(orders.StatusPaid),
			),
			expReq: testutils.NewTestOrder(
				testutils.WithStatus(orders.StatusPaid),
			),
			expErr: deterrs.NewDetErr(
				deterrs.OrderActionNotPermittedByStatus,
			),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.req.Close()
			testutils.AssertError(t, c.expErr, err)
			testutils.ValidateOrder(t, c.expReq, c.req)
		})
	}
}

func TestCancel(t *testing.T) {
	cases := []struct {
		name         string
		req          *orders.Order
		cancelReason string
		expReq       *orders.Order
		expErr       error
	}{
		{
			name: "Успешная отмена новой заявки",
			req: testutils.NewTestOrder(
				testutils.WithStatus(orders.StatusNew),
			),
			cancelReason: testutils.FilledCancelReason,
			expReq: testutils.NewTestOrder(
				testutils.WithStatus(orders.StatusCanceled),
				testutils.WithCancelReason(testutils.FilledCancelReason),
				testutils.WithScheduledFor(nil),
			),
			expErr: nil,
		},
		{
			name: "Успешная отмена запланированной заявки",
			req: testutils.NewTestOrder(
				testutils.WithStatus(orders.StatusScheduled),
			),
			cancelReason: testutils.FilledCancelReason,
			expReq: testutils.NewTestOrder(
				testutils.WithStatus(orders.StatusCanceled),
				testutils.WithCancelReason(testutils.FilledCancelReason),
				testutils.WithScheduledFor(nil),
			),
			expErr: nil,
		},
		{
			name: "Попытка отмены отменённой заявки",
			req: testutils.NewTestOrder(
				testutils.WithStatus(orders.StatusCanceled),
				testutils.WithCancelReason(testutils.FilledCancelReason),
			),
			cancelReason: "Другая причина отмены",
			expReq: testutils.NewTestOrder(
				testutils.WithStatus(orders.StatusCanceled),
				testutils.WithCancelReason(testutils.FilledCancelReason),
			),
			expErr: deterrs.NewDetErr(
				deterrs.OrderActionNotPermittedByStatus,
			),
		},
		{
			name: "Попытка отмены оплаченной заявки",
			req: testutils.NewTestOrder(
				testutils.WithStatus(orders.StatusPaid),
			),
			cancelReason: testutils.FilledCancelReason,
			expReq: testutils.NewTestOrder(
				testutils.WithStatus(orders.StatusPaid),
			),
			expErr: deterrs.NewDetErr(
				deterrs.OrderActionNotPermittedByStatus,
			),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.req.Cancel(c.cancelReason)
			testutils.AssertError(t, c.expErr, err)
			testutils.ValidateOrder(t, c.expReq, c.req)
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
		req           *orders.Order
		patchedFields *orders.OrderPatcher
		expReq        *orders.Order
		expErr        error
	}{
		{
			name: "Успешная модификация полей заявки",
			req:  testutils.NewTestOrder(),
			patchedFields: &orders.OrderPatcher{
				ClientName:          &patchedClientName,
				ClientPhone:         &patchedClientPhone,
				Address:             &patchedAddress,
				ClientDescription:   &patchedCliendDescription,
				EmployeeDescription: &patchedEmployeeDescription,
			},
			expReq: testutils.NewTestOrder(
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
			req: testutils.NewTestOrder(
				testutils.WithStatus(orders.StatusCanceled),
			),
			patchedFields: &orders.OrderPatcher{
				ClientName:          &patchedClientName,
				ClientPhone:         &patchedClientPhone,
				Address:             &patchedAddress,
				ClientDescription:   &patchedCliendDescription,
				EmployeeDescription: &patchedEmployeeDescription,
			},
			expReq: testutils.NewTestOrder(
				testutils.WithStatus(orders.StatusCanceled),
			),
			expErr: deterrs.NewDetErr(
				deterrs.OrderActionNotPermittedByStatus,
			),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.req.Patch(c.patchedFields)
			testutils.AssertError(t, c.expErr, err)
			testutils.ValidateOrder(t, c.expReq, c.req)
		})
	}
}
