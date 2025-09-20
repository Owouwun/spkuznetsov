package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Owouwun/spkuznetsov/internal/core/logic/orders"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// --- Mock service ---------------------------------------------------------

type MockOrderService struct {
	CreateFn      func(ctx context.Context, pord *orders.PrimaryOrder) (uuid.UUID, error)
	GetByIDFn     func(ctx context.Context, id uuid.UUID) (*orders.Order, error)
	GetAllFn      func(ctx context.Context) ([]*orders.Order, error)
	PrescheduleFn func(ctx context.Context, id uuid.UUID, scheduledFor *time.Time) error
	AssignFn      func(ctx context.Context, id uuid.UUID, empID uint) error
	ScheduleFn    func(ctx context.Context, id uuid.UUID, scheduledFor *time.Time) error
	ProgressFn    func(ctx context.Context, id uuid.UUID, empDescr string) error
	CompleteFn    func(ctx context.Context, id uuid.UUID) error
	CloseFn       func(ctx context.Context, id uuid.UUID) error
	CancelFn      func(ctx context.Context, id uuid.UUID, reason string) error
}

func (m *MockOrderService) Create(ctx context.Context, pord *orders.PrimaryOrder) (uuid.UUID, error) {
	if m.CreateFn == nil {
		return uuid.Nil, nil
	}
	return m.CreateFn(ctx, pord)
}
func (m *MockOrderService) GetByID(ctx context.Context, id uuid.UUID) (*orders.Order, error) {
	if m.GetByIDFn == nil {
		return nil, nil
	}
	return m.GetByIDFn(ctx, id)
}
func (m *MockOrderService) GetAll(ctx context.Context) ([]*orders.Order, error) {
	if m.GetAllFn == nil {
		return nil, nil
	}
	return m.GetAllFn(ctx)
}
func (m *MockOrderService) Preschedule(ctx context.Context, id uuid.UUID, scheduledFor *time.Time) error {
	if m.PrescheduleFn == nil {
		return nil
	}
	return m.PrescheduleFn(ctx, id, scheduledFor)
}
func (m *MockOrderService) Assign(ctx context.Context, id uuid.UUID, empID uint) error {
	if m.AssignFn == nil {
		return nil
	}
	return m.AssignFn(ctx, id, empID)
}
func (m *MockOrderService) Schedule(ctx context.Context, id uuid.UUID, scheduledFor *time.Time) error {
	if m.ScheduleFn == nil {
		return nil
	}
	return m.ScheduleFn(ctx, id, scheduledFor)
}
func (m *MockOrderService) Progress(ctx context.Context, id uuid.UUID, empDescr string) error {
	if m.ProgressFn == nil {
		return nil
	}
	return m.ProgressFn(ctx, id, empDescr)
}
func (m *MockOrderService) Complete(ctx context.Context, id uuid.UUID) error {
	if m.CompleteFn == nil {
		return nil
	}
	return m.CompleteFn(ctx, id)
}
func (m *MockOrderService) Close(ctx context.Context, id uuid.UUID) error {
	if m.CloseFn == nil {
		return nil
	}
	return m.CloseFn(ctx, id)
}
func (m *MockOrderService) Cancel(ctx context.Context, id uuid.UUID, reason string) error {
	if m.CancelFn == nil {
		return nil
	}
	return m.CancelFn(ctx, id, reason)
}

// --- Helpers --------------------------------------------------------------

func performRequest(handler http.Handler, method, path string, body []byte, contentType string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, bytes.NewReader(body))
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	return w
}

func requireJSONObj(t *testing.T, body []byte) {
	t.Helper()
	var m map[string]interface{}
	if err := json.Unmarshal(body, &m); err != nil {
		t.Fatalf("expected JSON object, unmarshal error: %v, body: %s", err, string(body))
	}
	if m == nil {
		t.Fatalf("expected JSON object, got nil (body: %s)", string(body))
	}
}

type MockSetupSimple func() *MockOrderService
type MockSetupWithCheck func() (*MockOrderService, func(t *testing.T))

// --- Tests ---------------

func TestGetAll_Handler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cases := []struct {
		name       string
		mockSetup  MockSetupSimple
		wantStatus int
	}{
		{
			name: "Успешно — возвращает список заявок",
			mockSetup: func() *MockOrderService {
				return &MockOrderService{
					GetAllFn: func(ctx context.Context) ([]*orders.Order, error) {
						return []*orders.Order{{}}, nil
					},
				}
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "Ошибка сервиса -> 500",
			mockSetup: func() *MockOrderService {
				return &MockOrderService{
					GetAllFn: func(ctx context.Context) ([]*orders.Order, error) {
						return nil, errors.New("boom")
					},
				}
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mock := tc.mockSetup()
			h := NewOrderHandler(mock)
			r := gin.New()
			r.GET("/orders", h.GetAll)

			w := performRequest(r, "GET", "/orders", nil, "")
			if w.Code != tc.wantStatus {
				t.Fatalf("want %d, got %d, body: %s", tc.wantStatus, w.Code, w.Body.String())
			}
		})
	}
}

func TestGetByID_Handler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testID := uuid.New()
	cases := []struct {
		name       string
		targetPath string
		mockSetup  MockSetupSimple
		wantStatus int
	}{
		{
			name:       "Неправильный UUID -> 400",
			targetPath: "/orders/not-a-uuid",
			mockSetup: func() *MockOrderService {
				return &MockOrderService{}
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Ошибка сервиса -> 500",
			targetPath: "/orders/" + testID.String(),
			mockSetup: func() *MockOrderService {
				return &MockOrderService{
					GetByIDFn: func(ctx context.Context, id uuid.UUID) (*orders.Order, error) {
						return nil, errors.New("db err")
					},
				}
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:       "Успех -> 200 и JSON-объект",
			targetPath: "/orders/" + testID.String(),
			mockSetup: func() *MockOrderService {
				return &MockOrderService{
					GetByIDFn: func(ctx context.Context, id uuid.UUID) (*orders.Order, error) {
						return &orders.Order{}, nil
					},
				}
			},
			wantStatus: http.StatusOK,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mock := tc.mockSetup()
			h := NewOrderHandler(mock)
			r := gin.New()
			r.GET("/orders/:id", h.GetByID)

			w := performRequest(r, "GET", tc.targetPath, nil, "")
			if w.Code != tc.wantStatus {
				t.Fatalf("want %d got %d body: %s", tc.wantStatus, w.Code, w.Body.String())
			}

			if tc.wantStatus == http.StatusOK {
				requireJSONObj(t, w.Body.Bytes())
			}
		})
	}
}

func TestCreate_Handler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	validJSON := []byte(`{}`)
	expectedID := uuid.New()

	cases := []struct {
		name       string
		body       []byte
		mockSetup  MockSetupSimple
		wantStatus int
		wantBody   string
	}{
		{
			name: "Некорректный JSON -> 400",
			body: []byte("not json"),
			mockSetup: func() *MockOrderService {
				return &MockOrderService{}
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Ошибка сервиса -> 500",
			body: validJSON,
			mockSetup: func() *MockOrderService {
				return &MockOrderService{
					CreateFn: func(ctx context.Context, pord *orders.PrimaryOrder) (uuid.UUID, error) {
						return uuid.Nil, errors.New("can't create")
					},
				}
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "Успешное создание -> 201 с id в теле",
			body: validJSON,
			mockSetup: func() *MockOrderService {
				return &MockOrderService{
					CreateFn: func(ctx context.Context, pord *orders.PrimaryOrder) (uuid.UUID, error) {
						return expectedID, nil
					},
				}
			},
			wantStatus: http.StatusCreated,
			wantBody:   expectedID.String(),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mock := tc.mockSetup()
			h := NewOrderHandler(mock)
			r := gin.New()
			r.POST("/orders", h.Create)

			w := performRequest(r, "POST", "/orders", tc.body, "application/json")
			if w.Code != tc.wantStatus {
				t.Fatalf("want %d got %d body: %s", tc.wantStatus, w.Code, w.Body.String())
			}
			if tc.wantBody != "" && !strings.Contains(w.Body.String(), tc.wantBody) {
				t.Fatalf("response body does not contain %q: %s", tc.wantBody, w.Body.String())
			}
		})
	}
}

func TestPreschedule_Handler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	id := uuid.New()
	validPayload, _ := json.Marshal(PrescheduleRequest{ScheduledFor: func() *time.Time { t := time.Now().Add(time.Hour); return &t }()})

	cases := []struct {
		name       string
		path       string
		body       []byte
		mockSetup  MockSetupSimple
		wantStatus int
	}{
		{
			name:       "Неверный UUID -> 400",
			path:       "/orders/bad/preschedule",
			body:       []byte(`{}`),
			mockSetup:  func() *MockOrderService { return &MockOrderService{} },
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Плохое тело -> 400",
			path:       "/orders/" + id.String() + "/preschedule",
			body:       []byte("notjson"),
			mockSetup:  func() *MockOrderService { return &MockOrderService{} },
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Сервис вернул validation error -> 400",
			path: "/orders/" + id.String() + "/preschedule",
			body: validPayload,
			mockSetup: func() *MockOrderService {
				return &MockOrderService{
					PrescheduleFn: func(ctx context.Context, id uuid.UUID, scheduledFor *time.Time) error {
						return errors.New("validation fail")
					},
				}
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Успех -> 200",
			path: "/orders/" + id.String() + "/preschedule",
			body: validPayload,
			mockSetup: func() *MockOrderService {
				return &MockOrderService{
					PrescheduleFn: func(ctx context.Context, id uuid.UUID, scheduledFor *time.Time) error {
						return nil
					},
				}
			},
			wantStatus: http.StatusOK,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mock := tc.mockSetup()
			h := NewOrderHandler(mock)
			r := gin.New()
			r.PATCH("/orders/:id/preschedule", h.Preschedule)

			w := performRequest(r, "PATCH", tc.path, tc.body, "application/json")
			if w.Code != tc.wantStatus {
				t.Fatalf("want %d got %d body: %s", tc.wantStatus, w.Code, w.Body.String())
			}
		})
	}
}

func TestAssign_Handler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	orderID := uuid.New()
	cases := []struct {
		name       string
		path       string
		mockSetup  MockSetupWithCheck
		wantStatus int
		wantEmp    uint
	}{
		{
			name:       "Неверный id заказа -> 400",
			path:       "/orders/bad/assign/1",
			mockSetup:  func() (*MockOrderService, func(t *testing.T)) { return &MockOrderService{}, nil },
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Неверный id сотрудника -> 400",
			path:       "/orders/" + orderID.String() + "/assign/notanint",
			mockSetup:  func() (*MockOrderService, func(t *testing.T)) { return &MockOrderService{}, nil },
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Сервис вернул ошибку -> 400",
			path: "/orders/" + orderID.String() + "/assign/7",
			mockSetup: func() (*MockOrderService, func(t *testing.T)) {
				return &MockOrderService{
					AssignFn: func(ctx context.Context, id uuid.UUID, empID uint) error {
						return errors.New("cannot assign")
					},
				}, nil
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Успешное присвоение -> 200 и empID передан в сервис",
			path: "/orders/" + orderID.String() + "/assign/42",
			mockSetup: func() (*MockOrderService, func(t *testing.T)) {
				var gotEmp uint
				mock := &MockOrderService{
					AssignFn: func(ctx context.Context, id uuid.UUID, empID uint) error {
						gotEmp = empID
						return nil
					},
				}
				return mock, func(t *testing.T) {
					if gotEmp != 42 {
						t.Fatalf("expected emp id 42, got %d", gotEmp)
					}
				}
			},
			wantStatus: http.StatusOK,
			wantEmp:    42,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mock, postCheck := tc.mockSetup()
			h := NewOrderHandler(mock)
			r := gin.New()
			r.PATCH("/orders/:id/assign/:empID", h.Assign)

			w := performRequest(r, "PATCH", tc.path, nil, "")
			if w.Code != tc.wantStatus {
				t.Fatalf("want %d got %d body: %s", tc.wantStatus, w.Code, w.Body.String())
			}
			if postCheck != nil {
				postCheck(t)
			}
		})
	}
}

func TestSchedule_Handler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	id := uuid.New()
	payload, _ := json.Marshal(PrescheduleRequest{ScheduledFor: func() *time.Time { t := time.Now().Add(2 * time.Hour); return &t }()})

	cases := []struct {
		name       string
		path       string
		body       []byte
		mockSetup  MockSetupSimple
		wantStatus int
	}{
		{
			name:       "Неверный UUID -> 400",
			path:       "/orders/bad/schedule",
			body:       []byte(`{}`),
			mockSetup:  func() *MockOrderService { return &MockOrderService{} },
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Плохое тело -> 400",
			path:       "/orders/" + id.String() + "/schedule",
			body:       []byte("nojson"),
			mockSetup:  func() *MockOrderService { return &MockOrderService{} },
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Сервис вернул ошибку -> 400",
			path: "/orders/" + id.String() + "/schedule",
			body: payload,
			mockSetup: func() *MockOrderService {
				return &MockOrderService{
					ScheduleFn: func(ctx context.Context, id uuid.UUID, scheduledFor *time.Time) error {
						return errors.New("can't schedule")
					},
				}
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Успех -> 200",
			path: "/orders/" + id.String() + "/schedule",
			body: payload,
			mockSetup: func() *MockOrderService {
				return &MockOrderService{
					ScheduleFn: func(ctx context.Context, id uuid.UUID, scheduledFor *time.Time) error { return nil },
				}
			},
			wantStatus: http.StatusOK,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mock := tc.mockSetup()
			h := NewOrderHandler(mock)
			r := gin.New()
			r.PATCH("/orders/:id/schedule", h.Schedule)

			w := performRequest(r, "PATCH", tc.path, tc.body, "application/json")
			if w.Code != tc.wantStatus {
				t.Fatalf("want %d got %d body: %s", tc.wantStatus, w.Code, w.Body.String())
			}
		})
	}
}

func TestProgress_Handler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	id := uuid.New()
	validBody, _ := json.Marshal(ProgressRequest{EmployeeDescription: "working"})

	cases := []struct {
		name       string
		path       string
		body       []byte
		mockSetup  MockSetupWithCheck
		wantStatus int
	}{
		{
			name:       "Неверный UUID -> 400",
			path:       "/orders/bad/progress",
			body:       []byte(`{}`),
			mockSetup:  func() (*MockOrderService, func(t *testing.T)) { return &MockOrderService{}, nil },
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Плохое тело -> 400",
			path:       "/orders/" + id.String() + "/progress",
			body:       []byte("nojson"),
			mockSetup:  func() (*MockOrderService, func(t *testing.T)) { return &MockOrderService{}, nil },
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Сервис вернул ошибку -> 400",
			path: "/orders/" + id.String() + "/progress",
			body: validBody,
			mockSetup: func() (*MockOrderService, func(t *testing.T)) {
				return &MockOrderService{
					ProgressFn: func(ctx context.Context, id uuid.UUID, empDescr string) error {
						return errors.New("cannot progress")
					},
				}, nil
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Успех -> 200",
			path: "/orders/" + id.String() + "/progress",
			body: validBody,
			mockSetup: func() (*MockOrderService, func(t *testing.T)) {
				return &MockOrderService{
					ProgressFn: func(ctx context.Context, id uuid.UUID, empDescr string) error { return nil },
				}, nil
			},
			wantStatus: http.StatusOK,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mock, _ := tc.mockSetup()
			h := NewOrderHandler(mock)
			r := gin.New()
			r.PATCH("/orders/:id/progress", h.Progress)

			w := performRequest(r, "PATCH", tc.path, tc.body, "application/json")
			if w.Code != tc.wantStatus {
				t.Fatalf("want %d got %d body: %s", tc.wantStatus, w.Code, w.Body.String())
			}
		})
	}
}

func TestComplete_Handler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	id := uuid.New()
	cases := []struct {
		name       string
		path       string
		mockSetup  MockSetupSimple
		wantStatus int
	}{
		{
			name:       "Неверный UUID -> 400",
			path:       "/orders/zzz/complete",
			mockSetup:  func() *MockOrderService { return &MockOrderService{} },
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Сервис вернул ошибку -> 400",
			path: "/orders/" + id.String() + "/complete",
			mockSetup: func() *MockOrderService {
				return &MockOrderService{
					CompleteFn: func(ctx context.Context, id uuid.UUID) error { return errors.New("bad") },
				}
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Успех -> 200",
			path: "/orders/" + id.String() + "/complete",
			mockSetup: func() *MockOrderService {
				return &MockOrderService{
					CompleteFn: func(ctx context.Context, id uuid.UUID) error { return nil },
				}
			},
			wantStatus: http.StatusOK,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mock := tc.mockSetup()
			h := NewOrderHandler(mock)
			r := gin.New()
			r.PATCH("/orders/:id/complete", h.Complete)

			w := performRequest(r, "PATCH", tc.path, nil, "")
			if w.Code != tc.wantStatus {
				t.Fatalf("want %d got %d body: %s", tc.wantStatus, w.Code, w.Body.String())
			}
		})
	}
}

func TestClose_Handler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	id := uuid.New()
	cases := []struct {
		name       string
		path       string
		mockSetup  MockSetupSimple
		wantStatus int
	}{
		{
			name:       "Неверный UUID -> 400",
			path:       "/orders/not-uuid/close",
			mockSetup:  func() *MockOrderService { return &MockOrderService{} },
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Сервис вернул ошибку -> 400",
			path: "/orders/" + id.String() + "/close",
			mockSetup: func() *MockOrderService {
				return &MockOrderService{
					CloseFn: func(ctx context.Context, id uuid.UUID) error { return errors.New("can't close") },
				}
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Успех -> 200",
			path: "/orders/" + id.String() + "/close",
			mockSetup: func() *MockOrderService {
				return &MockOrderService{
					CloseFn: func(ctx context.Context, id uuid.UUID) error { return nil },
				}
			},
			wantStatus: http.StatusOK,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mock := tc.mockSetup()
			h := NewOrderHandler(mock)
			r := gin.New()
			r.PATCH("/orders/:id/close", h.Close)

			w := performRequest(r, "PATCH", tc.path, nil, "")
			if w.Code != tc.wantStatus {
				t.Fatalf("want %d got %d body: %s", tc.wantStatus, w.Code, w.Body.String())
			}
		})
	}
}

func TestCancel_Handler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	id := uuid.New()
	validBody, _ := json.Marshal(CancelRequest{CancelReason: "no longer needed"})

	cases := []struct {
		name       string
		path       string
		body       []byte
		mockSetup  MockSetupWithCheck
		wantStatus int
	}{
		{
			name:       "Неверный UUID -> 400",
			path:       "/orders/bad/cancel",
			body:       []byte(`{}`),
			mockSetup:  func() (*MockOrderService, func(t *testing.T)) { return &MockOrderService{}, nil },
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Плохое тело -> 400",
			path:       "/orders/" + id.String() + "/cancel",
			body:       []byte("nojson"),
			mockSetup:  func() (*MockOrderService, func(t *testing.T)) { return &MockOrderService{}, nil },
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Сервис вернул ошибку -> 400",
			path: "/orders/" + id.String() + "/cancel",
			body: validBody,
			mockSetup: func() (*MockOrderService, func(t *testing.T)) {
				return &MockOrderService{CancelFn: func(ctx context.Context, id uuid.UUID, reason string) error { return errors.New("can't cancel") }}, nil
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Успех -> 200 и reason проброшен в сервис",
			path: "/orders/" + id.String() + "/cancel",
			body: validBody,
			mockSetup: func() (*MockOrderService, func(t *testing.T)) {
				var gotReason string
				mock := &MockOrderService{
					CancelFn: func(ctx context.Context, id uuid.UUID, reason string) error {
						gotReason = reason
						return nil
					},
				}
				return mock, func(t *testing.T) {
					if gotReason != "no longer needed" {
						t.Fatalf("expected reason %q, got %q", "no longer needed", gotReason)
					}
				}
			},
			wantStatus: http.StatusOK,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mock, postCheck := tc.mockSetup()
			h := NewOrderHandler(mock)
			r := gin.New()
			r.PATCH("/orders/:id/cancel", h.Cancel)

			w := performRequest(r, "PATCH", tc.path, tc.body, "application/json")
			if w.Code != tc.wantStatus {
				t.Fatalf("want %d got %d body: %s", tc.wantStatus, w.Code, w.Body.String())
			}
			if postCheck != nil {
				postCheck(t)
			}
		})
	}
}
