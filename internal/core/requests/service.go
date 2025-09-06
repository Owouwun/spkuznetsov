package requests

import (
	"time"

	"github.com/Owouwun/ipkuznetsov/internal/core/auth"
)

// TODO Решить между возвращением *Request и *RequestHistoryEntry

// Оформить новую заявку
// Ключи inputData: ClientName, ClientPhone, Address и, возможно, ClientDescription
func (preq *PrimaryRequest) CreateNewRequest() (*Request, error) {
	return nil, ErrNotImplemented
}

// Назначить предварительную дату работ
func (req *Request) Preschedule(date *time.Time) error {
	return ErrNotImplemented
}

// Назначить ответственного сотрудника
func (req *Request) Assign(emp *auth.Employee) error {
	return ErrNotImplemented
}

// Назначить точную дату выполнения работ
func (req *Request) Schedule(date *time.Time) error {
	return ErrNotImplemented
}

// Определить предварительную дату выполнения работ как точную
func (req *Request) ConfirmSchedule() error {
	return ErrNotImplemented
}

// Описать частично проведённые работы
func (req *Request) Progress(empDescription string) error {
	return ErrNotImplemented
}

// Пометить заявку как выполненную
func (req *Request) Complete() error {
	return ErrNotImplemented
}

// Закрыть заявку (после получения оплаты)
func (req *Request) Close() error {
	return ErrNotImplemented
}

// Отменить заявку с указанием причины
func (req *Request) Cancel(cause string) error {
	return ErrNotImplemented
}

// Модифицировать поля заявки
func (req *Request) Patch(patchedFields *RequestPatcher) error {
	return ErrNotImplemented
}
