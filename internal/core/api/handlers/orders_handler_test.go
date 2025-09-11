package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Owouwun/spkuznetsov/internal/core/api/handlers"
	"github.com/Owouwun/spkuznetsov/internal/core/api/mocks"
	"github.com/Owouwun/spkuznetsov/internal/core/logic/orders"
	"github.com/Owouwun/spkuznetsov/internal/testutils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateNewOrder_Success(t *testing.T) {
	const endpointLink = "/orders"

	// Задаём входные данные
	pord := orders.PrimaryOrder{
		ClientName:        testutils.ClientName,
		ClientPhone:       testutils.ClientPhone,
		Address:           testutils.Address,
		ClientDescription: testutils.ClientDescription,
	}
	jsonBytes, err := json.Marshal(pord)
	assert.NoError(t, err)

	// Подготавливаем тестовый запрос к эндпоинту
	req, err := http.NewRequest(http.MethodPost, endpointLink, bytes.NewBuffer(jsonBytes))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	// Задаём ожидаемый результат работы мока
	expectedOrder := testutils.NewTestOrder()

	// Создаём мок бизнес логики и задаём ему поведение (на входе — первичная заявка, на выходе — ожидаемая заявка)
	mockService := new(mocks.MockOrderService)
	mockService.On("CreateNewOrder", mock.AnythingOfType("*orders.PrimaryOrder")).Return(expectedOrder, nil).Once()

	// Создаём обработчик с созданным моком
	handler := handlers.NewOrderHandler(mockService)

	// Настраиваем gin под тестирование эндпоинта
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST(endpointLink, handler.CreateNewOrder)

	// Выполняем подготовленный запрос, сохранив ответ
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Проверяем результат работы запроса
	// Проверка статуса
	assert.Equal(t, http.StatusCreated, w.Code)
	// Проверка тела
	var response *handlers.CreateNewOrderResponseBody
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedOrder, response.Order)

	// Убеждаемся, что мок завершил работу
	mockService.AssertExpectations(t)
}
