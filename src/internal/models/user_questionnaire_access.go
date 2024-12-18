package models

import (
	"gorm.io/gorm"
	"time"
)

type UserQuestionnaireAccess struct {
	gorm.Model
	BaseValidator
	UserID          uint `gorm:"not null"`
	QuestionnaireID uint `gorm:"not null"`
	RoleID          uint `gorm:"not null"`
	ExpiresAt       *time.Time
	IsActive        bool `gorm:"default:true"`
	Role            *QuestionnaireRole
	TransactionID   *uint
	User            User `gorm:"foreignKey:UserID"`
}

func (q *UserQuestionnaireAccess) Validate() error {
	return q.BaseValidator.Validate(q)
}
