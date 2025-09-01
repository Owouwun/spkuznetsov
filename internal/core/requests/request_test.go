package requests

import (
	"fmt"
	"reflect"
	"testing"
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

func TestNewRequest(t *testing.T) {
	cases := []struct {
		name   string
		PReq   *PrimaryRequest
		expReq *Request
		err    error
	}{
		{
			name: "Успешное создание",
			PReq: &PrimaryRequest{
				ClientName:        "Test Client",
				ClientPhone:       "1234567890",
				Address:           "Test Address",
				ClientDescription: "Test Description",
			},
			expReq: &Request{
				ClientName:        "Test Client",
				ClientPhone:       "1234567890",
				Address:           "Test Address",
				ClientDescription: "Test Description",
				CancelReason:      "",
				Status:            StatusNew,
				ScheduledFor:      nil,
				ConfirmedSchedule: false,
				Done:              false,
				Paid:              false,
			},
			err: nil,
		},
		{
			name: "Слишком короткий телефон клиента",
			PReq: &PrimaryRequest{
				ClientName:        "Test Client",
				ClientPhone:       "123",
				Address:           "Test Address",
				ClientDescription: "Test Description",
			},
			expReq: nil,
			err:    ErrPhoneDigitsCount,
		},
		{
			name: "Телефон клиента содержит буквы",
			PReq: &PrimaryRequest{
				ClientName:        "Test Client",
				ClientPhone:       "Вот мой номер телефона: +71112223344",
				Address:           "Test Address",
				ClientDescription: "Test Description",
			},
			expReq: nil,
			err:    ErrPhoneWrongSymbols,
		},
		{
			name: "Правильный телефон с декораторами",
			PReq: &PrimaryRequest{
				ClientName:        "Test Client",
				ClientPhone:       "8(111)222-33-44 ",
				Address:           "Test Address",
				ClientDescription: "Test Description",
			},
			expReq: &Request{
				ClientName:        "Test Client",
				ClientPhone:       "+71112223344",
				Address:           "Test Address",
				ClientDescription: "Test Description",
				CancelReason:      "",
				Status:            StatusNew,
				ScheduledFor:      nil,
				ConfirmedSchedule: false,
				Done:              false,
				Paid:              false,
			},
			err: ErrPhoneWrongSymbols,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req, err := CreateNewRequest(c.PReq)
			if err != c.err {
				t.Errorf("Expected error: '%s', got: '%s'", c.err, err)
			}
			err = req.validate(*c.expReq)
			if err != nil {
				t.Error(err)
			}
		})
	}
}

func TestPreschedule(t *testing.T)

func TestAssign(t *testing.T)

func TestSchedule(t *testing.T)

func TestConfirmSchedule(t *testing.T)

func TestProgress(t *testing.T)

func TestComplete(t *testing.T)

func TestClose(t *testing.T)

func TestCancel(t *testing.T)

func TestModify(t *testing.T)
