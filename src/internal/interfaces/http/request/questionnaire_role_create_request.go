package request

import (
	"github.com/QBG-P2/Voting-System/internal/models"
	"gorm.io/gorm"
)

type QuestionnaireRoleCreateRequest struct {
	gorm.Model
	models.BaseValidator `json:"-"`
	Name                 string `json:"name,omitempty" validate:"required"`
	QuestionnaireID      uint   `json:"questionnaire_id,omitempty" validate:"required"`
}

func (q *QuestionnaireRoleCreateRequest) Validate() error {
	return q.BaseValidator.Validate(q)
}

func (q *QuestionnaireRoleCreateRequest) MapToModel(model interface{}) error {
	return q.BaseValidator.MapToModel(q, model)
}
