package request

import (
	"github.com/QBG-P2/Voting-System/internal/models"
)

type UpdateResponseRequest struct {
	models.BaseValidator `json:"-"`
	Content              string `json:"content" validate:"required,min=1,max=500"`
	AttachmentURL        string `json:"attachment_url" validate:"omitempty,url"`
}

func (v *UpdateResponseRequest) Validate() error {
	return v.BaseValidator.Validate(v)
}

func (v *UpdateResponseRequest) MapToModel(model interface{}) error {
	return v.BaseValidator.MapToModel(v, model)
}
