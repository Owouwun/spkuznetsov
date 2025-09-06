package requests

import (
	"time"

	"github.com/Owouwun/ipkuznetsov/internal/core/auth"
	"github.com/Owouwun/ipkuznetsov/internal/utils"
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
// Ключи inputData: ClientName, ClientPhone, Address и, возможно, ClientDescription
func (preq *PrimaryRequest) CreateNewRequest() (*Request, error) {
	if preq.ClientName == "" {
		return nil, NewErrEmptyField("ClientName")
	}
	if preq.ClientPhone == "" {
		return nil, NewErrEmptyField("ClientPhone")
	}
	if preq.Address == "" {
		return nil, NewErrEmptyField("Address")
	}

	req := &Request{
		ClientName:        preq.ClientName,
		ClientPhone:       preq.ClientPhone,
		Address:           preq.Address,
		ClientDescription: preq.ClientDescription,
		PublicLink:        utils.GenerateRandomString(),
		Status:            StatusNew,
	}
	return req, nil
}

// Назначить предварительную дату работ
func (req *Request) Preschedule(date *time.Time) error {
	validStatuses := []Status{
		StatusNew,
		StatusPrescheduled,
		StatusAssigned,
		StatusScheduled,
		StatusInProgress,
	}

	if !req.Status.isValid(&validStatuses) {
		return ErrActionNotPermittedByStatus
	}

	if date != nil && date.Before(time.Now()) {
		return ErrInvalidDate
	}

	req.Status = StatusPrescheduled
	req.ScheduledFor = date
	return nil
}

// Назначить ответственного сотрудника
func (req *Request) Assign(emp *auth.Employee) error {
	invalidStatuses := []Status{
		StatusNew,
		StatusCanceled,
	}

	if req.Status.isInvalid(&invalidStatuses) {
		return ErrActionNotPermittedByStatus
	}

	req.Employee = emp
	req.Status = StatusAssigned
	return nil
}

// Назначить точную дату выполнения работ
func (req *Request) Schedule(date *time.Time) error {
	validStatuses := []Status{
		StatusAssigned,
		StatusInProgress,
	}

	if date == nil {
		return NewErrEmptyField("date")
	}
	if date.Before(time.Now()) {
		return ErrInvalidDate
	}
	if !req.Status.isValid(&validStatuses) {
		return ErrActionNotPermittedByStatus
	}

	req.Status = StatusScheduled
	req.ScheduledFor = date
	return nil
}

// Определить предварительную дату выполнения работ как точную
func (req *Request) ConfirmSchedule() error {
	validStatuses := []Status{
		StatusAssigned,
	}

	if !req.Status.isValid(&validStatuses) {
		return ErrActionNotPermittedByStatus
	}

	if req.ScheduledFor == nil {
		return ErrInvalidDate
	}

	req.Status = StatusScheduled
	return nil
}

// Описать частично проведённые работы
func (req *Request) Progress(empDescription string) error {
	validStatuses := []Status{
		StatusScheduled,
	}

	if !req.Status.isValid(&validStatuses) {
		return ErrActionNotPermittedByStatus
	}

	req.Status = StatusInProgress
	req.ScheduledFor = nil
	req.EmployeeDescription = empDescription
	return nil
}

// Пометить заявку как выполненную
func (req *Request) Complete() error {
	validStatuses := []Status{
		StatusInProgress,
	}

	if !req.Status.isValid(&validStatuses) {
		return ErrActionNotPermittedByStatus
	}

	req.Status = StatusDone
	return nil
}

// Закрыть заявку (после получения оплаты)
func (req *Request) Close() error {
	validStatuses := []Status{
		StatusDone,
	}

	if !req.Status.isValid(&validStatuses) {
		return ErrActionNotPermittedByStatus
	}

	req.Status = StatusPaid
	return nil
}

// Отменить заявку с указанием причины
func (req *Request) Cancel(cause string) error {
	invalidStatuses := []Status{
		StatusPaid,
		StatusCanceled,
	}

	if req.Status.isInvalid(&invalidStatuses) {
		return ErrActionNotPermittedByStatus
	}

	req.CancelReason = &cause
	req.Status = StatusCanceled
	req.ScheduledFor = nil
	return nil
}

// Модифицировать поля заявки
func (req *Request) Patch(patchedFields *RequestPatcher) error {
	invalidStatuses := []Status{
		StatusPaid,
		StatusCanceled,
	}

	if req.Status.isInvalid(&invalidStatuses) {
		return ErrActionNotPermittedByStatus
	}

	if patchedFields.ClientName != nil {
		req.ClientName = *patchedFields.ClientName
	}
	if patchedFields.ClientPhone != nil {
		req.ClientPhone = *patchedFields.ClientPhone
	}
	if patchedFields.Address != nil {
		req.Address = *patchedFields.Address
	}
	if patchedFields.ClientDescription != nil {
		req.ClientDescription = *patchedFields.ClientDescription
	}
	if patchedFields.EmployeeDescription != nil {
		req.EmployeeDescription = *patchedFields.EmployeeDescription
	}
	return nil
}
