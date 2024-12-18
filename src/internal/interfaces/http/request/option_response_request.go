package request

import "github.com/QBG-P2/Voting-System/internal/models"

type GetOptionsRequest struct {
	models.BaseValidator `json:"-"`
	QuestionID           uint              `query:"question_id" validate:"required"`
	Pagination           models.Pagination `query:"page,page_size"`
}

func (r *GetOptionsRequest) Validate() error {
	return r.BaseValidator.Validate(r)
}

func (r *GetOptionsRequest) MapToModel(model interface{}) error {
	return r.BaseValidator.MapToModel(r, model)
}
