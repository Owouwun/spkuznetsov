package entities

import "github.com/Owouwun/spkuznetsov/internal/core/logic/auth"

type EmployeeEntity struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"not null"`
}

func (EmployeeEntity) TableName() string {
	return "public.employees"
}

func (ee *EmployeeEntity) ToLogicEmployee() *auth.Employee {
	if ee == nil {
		return nil
	}
	return &auth.Employee{
		ID:   ee.ID,
		Name: ee.Name,
	}
}

func NewEmployeeEntityFromLogic(emp *auth.Employee) *EmployeeEntity {
	if emp == nil {
		return nil
	}

	ee := &EmployeeEntity{
		ID:   emp.ID,
		Name: emp.Name,
	}

	return ee
}
