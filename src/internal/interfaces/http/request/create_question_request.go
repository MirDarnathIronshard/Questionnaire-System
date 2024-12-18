package request

import "github.com/QBG-P2/Voting-System/internal/models"

type CreateQuestionRequest struct {
	models.BaseValidator `json:"-"`
	Text                 string              `json:"text" validate:"required,min=1,max=500"`
	Type                 models.QuestionType `json:"type" validate:"required,oneof=MultipleChoice ShortAnswer Decriptive"`
	IsConditional        bool                `json:"is_conditional"`
	QuestionnaireID      uint                `json:"questionnaire_id" validate:"required"`
	Order                int                 `json:"order" validate:"gte=1"`
	Condition            *string             `json:"condition" validate:"omitempty"`
	MediaURL             *string             `json:"media_url" validate:"omitempty,url"`
	CorrectAnswer        *string             `json:"correct_answer" validate:"omitempty"`
	Options              []OptionRequest     `json:"options,omitempty" validate:"dive"`
}

type OptionRequest struct {
	Text string `json:"text" validate:"required,min=1,max=255"`
}

func (v *CreateQuestionRequest) Validate() error {
	return v.BaseValidator.Validate(v)
}
func (v *CreateQuestionRequest) MapToModel(model interface{}) error {
	return v.BaseValidator.MapToModel(v, model)
}
