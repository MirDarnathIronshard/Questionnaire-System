package request

import (
	"github.com/QBG-P2/Voting-System/internal/models"
	"time"
)

type UpdateProfileRequest struct {
	models.BaseValidator `json:"-"`
	Email                string    `json:"email" validate:"omitempty,email"`
	FirstName            string    `json:"first_name" validate:"omitempty,min=2,max=50"`
	LastName             string    `json:"last_name" validate:"omitempty,min=2,max=50"`
	BirthDate            time.Time `json:"birth_date" validate:"omitempty"`
	City                 string    `json:"city" validate:"omitempty,min=2,max=50"`
	Gender               string    `json:"gender" validate:"omitempty,oneof=male female other"`
}

func (v *UpdateProfileRequest) Validate() error {
	return v.BaseValidator.Validate(v)
}

func (v *UpdateProfileRequest) MapToModel(model interface{}) error {
	return v.BaseValidator.MapToModel(v, model)
}
