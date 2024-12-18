package request

import (
	"github.com/QBG-P2/Voting-System/internal/models"
	"time"
)

type QuestionnaireUpdateRequest struct {
	models.BaseValidator `json:"-"`
	Title                string                    `json:"title"`
	Description          string                    `json:"description"`
	StartTime            time.Time                 `json:"start_time"`
	EndTime              time.Time                 `json:"end_time"`
	StepType             models.StepType           `json:"step_type"`
	AllowBacktrack       bool                      `json:"allow_backtrack"`
	AllowedGenders       models.AllowedGenders     `json:"allowed_genders" validate:"required,oneof=male female all"`
	MaxAttempts          int                       `json:"max_attempts"`
	AnonymityLevel       models.TypeAnonymityLevel `json:"anonymity_level"`
	ResponseEditDeadline time.Time                 `json:"response_edit_deadline"`
	OwnerID              uint                      `json:"owner_id"`
}

func (q *QuestionnaireUpdateRequest) Validate() error {
	return q.BaseValidator.Validate(q)
}

func (q *QuestionnaireUpdateRequest) MapToModel(model interface{}) error {
	return q.BaseValidator.MapToModel(q, model)
}
