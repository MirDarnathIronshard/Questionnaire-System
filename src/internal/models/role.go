package models

import "gorm.io/gorm"

type Role struct {
	gorm.Model
	BaseValidator
	Name        string       `gorm:"uniqueIndex;not null" json:"name"`
	Permissions []Permission `gorm:"many2many:role_permissions;" json:"permissions"`
	Users       []User       `gorm:"many2many:user_roles;" json:"-"`
}

func (q *Role) Validate() error {
	return q.BaseValidator.Validate(q)
}
