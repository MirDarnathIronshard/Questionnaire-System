package request

import "github.com/QBG-P2/Voting-System/internal/models"

type AssignRoleRequest struct {
	models.BaseValidator `json:"-"`
	QuestionnaireID      uint   `json:"questionnaire_id" validate:"required"`
	UserID               uint   `json:"user_id" validate:"required"`
	Role                 string `json:"role" validate:"required,oneof=admin editor viewer"`
}

func (v *AssignRoleRequest) Validate() error {
	return v.BaseValidator.Validate(v)
}

func (v *AssignRoleRequest) MapToModel(model interface{}) error {
	return v.BaseValidator.MapToModel(v, model)
}
