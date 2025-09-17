package repository_auth

import (
	"context"

	"github.com/Owouwun/spkuznetsov/internal/core/logic/auth"
	"github.com/Owouwun/spkuznetsov/internal/core/repository/entities"
	"gorm.io/gorm"
)

type GormEmployeeRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) auth.AuthRepository {
	return &GormEmployeeRepository{db: db}
}

func (r *GormEmployeeRepository) getEntityByID(ctx context.Context, id uint) (*entities.EmployeeEntity, error) {
	var employeeEntity *entities.EmployeeEntity
	result := r.db.WithContext(ctx).
		First(&employeeEntity, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}

	return employeeEntity, nil
}

func (r *GormEmployeeRepository) CreateEmployee(ctx context.Context, emp *auth.Employee) (uint, error) {
	orderEntity := entities.NewEmployeeEntityFromLogic(emp)

	result := r.db.WithContext(ctx).Create(&orderEntity)
	if result.Error != nil {
		return 0, result.Error
	}

	return orderEntity.ID, nil
}

func (r *GormEmployeeRepository) GetEmployeeByID(ctx context.Context, id uint) (*auth.Employee, error) {
	employeeEntity, err := r.getEntityByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return employeeEntity.ToLogicEmployee(), nil
}

func (r *GormEmployeeRepository) GetEmployees(ctx context.Context) ([]*auth.Employee, error) {
	var employeeEntities []entities.EmployeeEntity
	result := r.db.WithContext(ctx).
		Find(&employeeEntities)

	if result.Error != nil {
		return nil, result.Error
	}

	var logicEmployees []*auth.Employee
	for _, entity := range employeeEntities {
		logicEmployees = append(logicEmployees, entity.ToLogicEmployee())
	}

	return logicEmployees, nil
}
