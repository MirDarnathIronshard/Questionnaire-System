package request

import (
	"github.com/QBG-P2/Voting-System/internal/models"
)

type GetResponsesRequest struct {
	models.BaseValidator `json:"-"`
	QuestionnaireID      uint `query:"questionnaire_id" validate:"required"`
	Page                 int  `query:"page" default:"1" json:"page"`
	PageSize             int  `query:"page_size" default:"20" json:"page_size"`
}

func (v *GetResponsesRequest) Validate() error {
	return v.BaseValidator.Validate(v)
}

func (v *GetResponsesRequest) MapToModel(model interface{}) error {
	return v.BaseValidator.MapToModel(v, model)
}
