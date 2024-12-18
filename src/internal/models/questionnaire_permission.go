package models

import "gorm.io/gorm"

type QuestionnairePermission struct {
	gorm.Model
	BaseValidator
	Name        string `gorm:"unique;not null"`
	Description string
	Action      string `gorm:"not null"` // e.g., "view", "edit", "delete", "manage_roles"
	Resource    string `gorm:"not null"` // e.g., "questionnaire", "responses", "results"
}

func (q *QuestionnairePermission) Validate() error {
	return q.BaseValidator.Validate(q)
}
