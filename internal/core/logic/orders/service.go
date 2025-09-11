package orders

import (
	"time"

	"github.com/Owouwun/spkuznetsov/internal/core/logic/auth"
	deterrs "github.com/Owouwun/spkuznetsov/internal/errors"
	"github.com/Owouwun/spkuznetsov/pkg/utils"
)

func (s *Status) isValid(validStatuses *[]Status) bool {
	for _, vs := range *validStatuses {
		if *s == vs {
			return true
		}
	}
	return false
}
func (s *Status) isInvalid(invalidStatuses *[]Status) bool {
	return s.isValid(invalidStatuses)
}

// Оформить новую заявку
func (pord *PrimaryOrder) CreateNewOrder() (*Order, error) {
	if pord.ClientName == "" {
		return nil, deterrs.NewDetErr(
			deterrs.EmptyField,
			deterrs.WithField("client name"),
		)
	}
	if pord.ClientPhone == "" {
		return nil, deterrs.NewDetErr(
			deterrs.EmptyField,
			deterrs.WithField("client phone"),
		)
	}
	if pord.Address == "" {
		return nil, deterrs.NewDetErr(
			deterrs.EmptyField,
			deterrs.WithField("address"),
		)
	}

	stdPN, err := utils.StandartizePhoneNumber(pord.ClientPhone)
	if err != nil {
		return nil, deterrs.NewDetErr(
			deterrs.InvalidValue,
			deterrs.WithField("client phone"),
			deterrs.WithOriginalError(err),
		)
	}

	ord := &Order{
		ClientName:        pord.ClientName,
		ClientPhone:       stdPN,
		Address:           pord.Address,
		ClientDescription: pord.ClientDescription,
		PublicLink:        utils.GenerateRandomString(),
		Status:            StatusNew,
	}
	return ord, nil
}

// Назначить предварительную дату работ
func (ord *Order) Preschedule(date *time.Time) error {
	validStatuses := []Status{
		StatusNew,
		StatusPrescheduled,
		StatusAssigned,
		StatusScheduled,
		StatusInProgress,
	}

	if !ord.Status.isValid(&validStatuses) {
		return deterrs.NewDetErr(
			deterrs.OrderActionNotPermittedByStatus,
		)
	}

	if date != nil {
		if err := utils.MustNotPast(date); err != nil {
			return deterrs.NewDetErr(
				deterrs.InvalidValue,
				deterrs.WithField("scheduled date"),
				deterrs.WithOriginalError(err),
			)
		}
	}

	ord.Status = StatusPrescheduled
	ord.ScheduledFor = date
	return nil
}

// Назначить ответственного сотрудника
func (ord *Order) Assign(emp *auth.Employee) error {
	invalidStatuses := []Status{
		StatusNew,
		StatusCanceled,
	}

	if ord.Status.isInvalid(&invalidStatuses) {
		return deterrs.NewDetErr(
			deterrs.OrderActionNotPermittedByStatus,
		)
	}

	ord.Employee = emp
	ord.Status = StatusAssigned
	return nil
}

// Назначить точную дату выполнения работ
func (ord *Order) Schedule(date *time.Time) error {
	validStatuses := []Status{
		StatusAssigned,
		StatusInProgress,
	}

	if date == nil {
		return deterrs.NewDetErr(
			deterrs.EmptyField,
			deterrs.WithField("scheduled date"),
		)
	}
	if err := utils.MustNotPast(date); err != nil {
		return deterrs.NewDetErr(
			deterrs.InvalidValue,
			deterrs.WithField("sheduled date"),
			deterrs.WithOriginalError(err),
		)
	}
	if !ord.Status.isValid(&validStatuses) {
		return deterrs.NewDetErr(
			deterrs.OrderActionNotPermittedByStatus,
		)
	}

	ord.Status = StatusScheduled
	ord.ScheduledFor = date
	return nil
}

// Определить предварительную дату выполнения работ как точную
func (ord *Order) ConfirmSchedule() error {
	validStatuses := []Status{
		StatusAssigned,
	}

	if !ord.Status.isValid(&validStatuses) {
		return deterrs.NewDetErr(
			deterrs.OrderActionNotPermittedByStatus,
		)
	}

	if ord.ScheduledFor == nil {
		return deterrs.NewDetErr(
			deterrs.EmptyField,
			deterrs.WithField("scheduled date"),
		)
	}

	ord.Status = StatusScheduled
	return nil
}

// Описать частично проведённые работы
func (ord *Order) Progress(empDescription string) error {
	validStatuses := []Status{
		StatusScheduled,
	}

	if !ord.Status.isValid(&validStatuses) {
		return deterrs.NewDetErr(
			deterrs.OrderActionNotPermittedByStatus,
		)
	}

	ord.Status = StatusInProgress
	ord.ScheduledFor = nil
	ord.EmployeeDescription = empDescription
	return nil
}

// Пометить заявку как выполненную
func (ord *Order) Complete() error {
	validStatuses := []Status{
		StatusInProgress,
	}

	if !ord.Status.isValid(&validStatuses) {
		return deterrs.NewDetErr(
			deterrs.OrderActionNotPermittedByStatus,
		)
	}

	ord.Status = StatusDone
	return nil
}

// Закрыть заявку (после получения оплаты)
func (ord *Order) Close() error {
	validStatuses := []Status{
		StatusDone,
	}

	if !ord.Status.isValid(&validStatuses) {
		return deterrs.NewDetErr(
			deterrs.OrderActionNotPermittedByStatus,
		)
	}

	ord.Status = StatusPaid
	return nil
}

// Отменить заявку с указанием причины
func (ord *Order) Cancel(cause string) error {
	invalidStatuses := []Status{
		StatusPaid,
		StatusCanceled,
	}

	if ord.Status.isInvalid(&invalidStatuses) {
		return deterrs.NewDetErr(
			deterrs.OrderActionNotPermittedByStatus,
		)
	}

	ord.CancelReason = &cause
	ord.Status = StatusCanceled
	ord.ScheduledFor = nil
	return nil
}

// Модифицировать поля заявки
func (ord *Order) Patch(patchedFields *OrderPatcher) error {
	invalidStatuses := []Status{
		StatusPaid,
		StatusCanceled,
	}

	if ord.Status.isInvalid(&invalidStatuses) {
		return deterrs.NewDetErr(
			deterrs.OrderActionNotPermittedByStatus,
		)
	}

	if patchedFields.ClientName != nil {
		ord.ClientName = *patchedFields.ClientName
	}
	if patchedFields.ClientPhone != nil {
		ord.ClientPhone = *patchedFields.ClientPhone
	}
	if patchedFields.Address != nil {
		ord.Address = *patchedFields.Address
	}
	if patchedFields.ClientDescription != nil {
		ord.ClientDescription = *patchedFields.ClientDescription
	}
	if patchedFields.EmployeeDescription != nil {
		ord.EmployeeDescription = *patchedFields.EmployeeDescription
	}
	return nil
}
