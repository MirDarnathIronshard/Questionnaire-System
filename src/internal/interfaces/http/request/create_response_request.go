package request

import "github.com/QBG-P2/Voting-System/internal/models"

type CreateResponseRequest struct {
	models.BaseValidator `json:"-"`
	QuestionnaireID      uint   `json:"questionnaire_id" validate:"required"`
	QuestionID           uint   `json:"question_id" validate:"required"`
	OptionID             uint   `json:"option_id" validate:"required"`
	Answer               string `json:"answer" validate:"required"`
	IsFinalized          bool   `json:"is_finalized" validate:"required"`
}

func (v *CreateResponseRequest) Validate() error {
	return v.BaseValidator.Validate(v)
}

func (v *CreateResponseRequest) MapToModel(model interface{}) error {
	return v.BaseValidator.MapToModel(v, model)
}
