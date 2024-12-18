package request

import "github.com/QBG-P2/Voting-System/internal/models"

type CreateOptionRequest struct {
	models.BaseValidator `json:"-"`
	Text                 string `json:"text" validate:"required,min=1,max=255"`
	QuestionID           uint   `json:"question_id" validate:"required"`
	QuestionnaireID      uint   `json:"questionnaire_id" validate:"required"`
}

func (r *CreateOptionRequest) Validate() error {
	return r.BaseValidator.Validate(r)
}

func (r *CreateOptionRequest) MapToModel(model interface{}) error {
	return r.BaseValidator.MapToModel(r, model)
}
