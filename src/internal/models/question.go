package models

import "gorm.io/gorm"

type QuestionType string

const (
	QuestionTypeMultipleChoice QuestionType = "MultipleChoice"
	QuestionTypeShortAnswer    QuestionType = "ShortAnswer"
	QuestionTypeDecriptive     QuestionType = "Decriptive"
)

type Question struct {
	gorm.Model
	BaseValidator
	Text            string `gorm:"not null"`
	Type            QuestionType
	IsConditional   bool `gorm:"default: false"`
	QuestionnaireID uint
	Order           int
	Condition       *string
	MediaURL        *string
	CorrectAnswer   *string
	Option          []Option   `gorm:"foreignkey:QuestionID"`
	Response        []Response `gorm:"foreignkey:QuestionID"`
}

func (q *Question) Validate() error {
	return q.BaseValidator.Validate(q)
}
