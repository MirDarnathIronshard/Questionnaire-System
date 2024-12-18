package request

import "github.com/QBG-P2/Voting-System/internal/models"

type UserCreateRequest struct {
	models.BaseValidator `json:"-"`
	NationalID           string `json:"national_id" validate:"required"`
	Email                string `json:"email" validate:"required"`
	Password             string `json:"password" validate:"required"`
	Gender               string `json:"gender" validate:"required"`
}
type GetOtpRequest struct {
	Email string `json:"mobileNumber" binding:"required,mobile,min=11,max=30"`
}

func (q *UserCreateRequest) Validate() error {
	return q.BaseValidator.Validate(q)
}

func (q *UserCreateRequest) MapToModel(model interface{}) error {
	return q.BaseValidator.MapToModel(q, model)
}
