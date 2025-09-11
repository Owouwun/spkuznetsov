package auth

import "gorm.io/gorm"

type Employee struct {
	Name string `json:"name"`
	// Gorm
	gorm.Model
}
