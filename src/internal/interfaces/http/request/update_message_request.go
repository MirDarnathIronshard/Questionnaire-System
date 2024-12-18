package request

import "github.com/QBG-P2/Voting-System/internal/models"

type UpdateMessageRequest struct {
	models.BaseValidator `json:"-"`
	Content              string `json:"content" validate:"required,min=1,max=1000"`
	AttachmentURL        string `json:"attachment_url" validate:"omitempty,url"`
}

func (q *UpdateMessageRequest) Validate() error {
	return q.BaseValidator.Validate(q)
}
func (q *UpdateMessageRequest) MapToModel(model interface{}) error {
	return q.BaseValidator.MapToModel(q, model)
}
