package request

import (
	"github.com/QBG-P2/Voting-System/internal/models"
	"gorm.io/gorm"
)

type QuestionnaireRoleUpdateRequest struct {
	gorm.Model           `json:"gorm_._model" validate:"required"`
	models.BaseValidator `json:"-" `
	ID                   uint   `json:"id,omitempty" validate:"required"`
	Name                 string `json:"name,omitempty" validate:"required"`
	QuestionnaireID      uint   `json:"questionnaire_id,omitempty" validate:"required"`
	UserID               uint   `json:"user_id,omitempty" validate:"required"`
}

func (q *QuestionnaireRoleUpdateRequest) Validate() error {
	return q.BaseValidator.Validate(q)
}

func (q *QuestionnaireRoleUpdateRequest) MapToModel(model interface{}) error {
	return q.BaseValidator.MapToModel(q, model)
}
