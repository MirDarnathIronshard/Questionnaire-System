package models

import "gorm.io/gorm"

type QuestionnaireRolePermission struct {
	gorm.Model
	BaseValidator
	QuestionnaireRoleID       uint                    `gorm:"primaryKey; autoIncrement:false"`
	QuestionnairePermissionID uint                    `gorm:"primaryKey; autoIncrement:false"`
	QuestionnaireRole         QuestionnaireRole       `gorm:"foreignKey:QuestionnaireRoleID"`
	QuestionnairePermission   QuestionnairePermission `gorm:"foreignKey:QuestionnairePermissionID"`
}

func (q *QuestionnaireRolePermission) Validate() error {
	return q.BaseValidator.Validate(q)
}
