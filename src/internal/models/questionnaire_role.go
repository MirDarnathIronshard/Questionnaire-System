package models

import (
	"gorm.io/gorm"
	"time"
)

type QuestionnaireRole struct {
	gorm.Model
	BaseValidator
	Name            string
	UserID          uint                      `gorm:"not null"`
	QuestionnaireID uint                      `gorm:"not null"`
	Permissions     []QuestionnairePermission `gorm:"many2many:questionnaire_role_permissions;"`
	ExpiresAt       *time.Time
	IsActive        bool `gorm:"default:true"`
}

func (q *QuestionnaireRole) Validate() error {
	return q.BaseValidator.Validate(q)
}

func (q *QuestionnaireRole) MapToModel(model interface{}) error {
	return q.BaseValidator.MapToModel(q, model)
}
