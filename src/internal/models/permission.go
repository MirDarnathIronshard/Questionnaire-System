package models

import "gorm.io/gorm"

type Permission struct {
	gorm.Model
	BaseValidator
	Name   string `gorm:"uniqueIndex;not null" json:"name"`
	Path   string `gorm:"not null" json:"path"`
	Method string `gorm:"not null" json:"method"`
	Roles  []Role `gorm:"many2many:role_permissions;" json:"-"`
}

func (q *Permission) Validate() error {
	return q.BaseValidator.Validate(q)
}
