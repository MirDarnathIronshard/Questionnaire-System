package request

import (
	"github.com/QBG-P2/Voting-System/internal/models"
)

type UpdateQuestionRequest struct {
	models.BaseValidator `json:"-"`
	Text                 string                `json:"text" validate:"required,min=1,max=500"`
	Type                 models.QuestionType   `json:"type" validate:"required,oneof=MultipleChoice ShortAnswer Decriptive"`
	IsConditional        bool                  `json:"is_conditional"`
	Order                int                   `json:"order" validate:"gte=1"`
	Condition            *string               `json:"condition" validate:"omitempty"`
	MediaURL             *string               `json:"media_url" validate:"omitempty,url"`
	CorrectAnswer        *string               `json:"correct_answer" validate:"omitempty"`
	Options              []OptionRequestUpdate `json:"options,omitempty" validate:"dive"`
}

type OptionRequestUpdate struct {
	ID   uint   `json:"id,omitempty"`
	Text string `json:"text" validate:"required,min=1,max=255"`
}

func (v *UpdateQuestionRequest) Validate() error {
	return v.BaseValidator.Validate(v)
}

func (v *UpdateQuestionRequest) MapToModel(model interface{}) error {
	return v.BaseValidator.MapToModel(v, model)
}
