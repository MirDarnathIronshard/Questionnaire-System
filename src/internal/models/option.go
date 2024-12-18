package models

import "gorm.io/gorm"

type Option struct {
	gorm.Model
	BaseValidator
	Text       string `gorm:"not null"`
	QuestionID uint
}

func (q *Option) Validate() error {
	return q.BaseValidator.Validate(q)
}
