package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/Owouwun/spkuznetsov/internal/core/logic/auth"
	"github.com/gin-gonic/gin"
)

// Задаёт методы бизнес-логики
type AuthService interface {
	CreateEmployee(ctx context.Context, name string) (uint, error)
	GetEmployeeByID(ctx context.Context, id uint) (*auth.Employee, error)
	GetEmployees(ctx context.Context) ([]*auth.Employee, error)
}

type AuthHandler struct {
	authService AuthService
}

func NewAuthHandler(as AuthService) *AuthHandler {
	return &AuthHandler{
		authService: as,
	}
}

func (h *AuthHandler) GetEmployees(c *gin.Context) {
	orders, err := h.authService.GetEmployees(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get employees", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)
}

func (h *AuthHandler) GetEmployeeByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid employee id", "details": err.Error()})
		return
	}

	employee, err := h.authService.GetEmployeeByID(c, uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get order", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, employee)
}

type NewEmployeeRequest struct {
	Name string `json:"name"`
}

func (h *AuthHandler) Create(c *gin.Context) {
	var newEmployee NewEmployeeRequest

	if err := c.ShouldBindJSON(&newEmployee); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body", "details": err.Error()})
		return
	}

	newEmployeeID, err := h.authService.CreateEmployee(c, newEmployee.Name)
	if err != nil {
		// TODO: Check on various errors
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create new employee", "details": err.Error()})
		return
	}

	// Отправляем успешный ответ с созданным заказом и статусом 201 Created.
	c.JSON(http.StatusCreated, newEmployeeID)
}
