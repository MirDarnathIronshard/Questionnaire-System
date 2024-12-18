package request

import (
	"github.com/QBG-P2/Voting-System/internal/models"
)

type UpdateRoleRequest struct {
	models.BaseValidator `json:"-"`
	Role                 string `json:"role" validate:"required,oneof=admin editor viewer"`
}

func (v *UpdateRoleRequest) Validate() error {
	return v.BaseValidator.Validate(v)
}

func (v *UpdateRoleRequest) MapToModel(model interface{}) error {
	return v.BaseValidator.MapToModel(v, model)
}
