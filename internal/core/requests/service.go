package requests

import (
	"time"

	"github.com/Owouwun/ipkuznetsov/internal/core/auth"
)

// TODO Решить между возвращением *Request и *RequestHistoryEntry

// Оформить новую заявку
// Ключи inputData: ClientName, ClientPhone, Address и, возможно, ClientDescription
func CreateNewRequest(preq *PrimaryRequest) (*Request, error)

// Назначить предварительную дату работ
func (req *Request) Preschedule(date *time.Time) error

// Назначить ответственного сотрудника
func (req *Request) Assign(emp *auth.Employee) error

// Назначить точную дату выполнения работ
func (req *Request) Schedule(date *time.Time) error

// Определить предварительную дату выполнения работ как точную
func (req *Request) ConfirmSchedule() error

// Описать частично проведённые работы
func (req *Request) Progress(empDescription string) error

// Пометить заявку как выполненную
func (req *Request) Complete() error

// Закрыть заявку (после получения оплаты)
func (req *Request) Close() error

// Отменить заявку с указанием причины
func (req *Request) Cancel(cause string) error

// Модифицировать поля заявки
func (req *Request) Modify(inputData map[string]string) error
