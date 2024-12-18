package request

import "github.com/QBG-P2/Voting-System/internal/models"

type GetQuestionsRequest struct {
	models.BaseValidator `json:"-"`
	QuestionnaireID      uint              `query:"questionnaire_id" validate:"required"`
	Pagination           models.Pagination `query:"page,default=1" validate:"required"`
}

func (v *GetQuestionsRequest) Validate() error {
	return v.BaseValidator.Validate(v)
}

func (v *GetQuestionsRequest) MapToModel(model interface{}) error {
	return v.BaseValidator.MapToModel(v, model)
}
