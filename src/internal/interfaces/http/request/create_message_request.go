package request

import (
	"github.com/QBG-P2/Voting-System/internal/models"
)

type CreateMessageRequest struct {
	models.BaseValidator `json:"-"`
	ChatID               uint   `json:"chat_id" validate:"required"`
	Content              string `json:"content" validate:"required,min=1,max=1000"`
	AttachmentURL        string `json:"attachment_url" validate:"omitempty,url"`
}

func (v *CreateMessageRequest) Validate() error {
	if err := v.BaseValidator.Validate(v); err != nil {
		return err
	}
	return nil
}

func (v *CreateMessageRequest) MapToModel(model interface{}) error {
	return v.BaseValidator.MapToModel(v, model)
}
