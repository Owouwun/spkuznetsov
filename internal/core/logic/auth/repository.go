package auth

import (
	"context"
)

type AuthRepository interface {
	GetEmployees(ctx context.Context) ([]*Employee, error)
	GetEmployeeByID(ctx context.Context, id uint) (*Employee, error)
	CreateEmployee(ctx context.Context, emp *Employee) (uint, error)
}

type AuthService struct {
	repo AuthRepository
}

func NewAuthService(repo AuthRepository) *AuthService {
	return &AuthService{
		repo: repo,
	}
}

func (s *AuthService) CreateEmployee(ctx context.Context, name string) (uint, error) {
	id, err := s.repo.CreateEmployee(ctx, &Employee{Name: name})
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *AuthService) GetEmployeeByID(ctx context.Context, id uint) (*Employee, error) {
	return s.repo.GetEmployeeByID(ctx, id)
}

func (s *AuthService) GetEmployees(ctx context.Context) ([]*Employee, error) {
	return s.repo.GetEmployees(ctx)
}
