package request

import "github.com/QBG-P2/Voting-System/internal/models"

type UpdateOptionRequest struct {
	models.BaseValidator `json:"-"`
	Text                 string `json:"text" validate:"required,min=1,max=255"`
	QuestionnaireID      uint   `json:"questionnaire_id" validate:"required"`
}

func (r *UpdateOptionRequest) Validate() error {
	return r.BaseValidator.Validate(r)
}

func (r *UpdateOptionRequest) MapToModel(model interface{}) error {
	return r.BaseValidator.MapToModel(r, model)
}
