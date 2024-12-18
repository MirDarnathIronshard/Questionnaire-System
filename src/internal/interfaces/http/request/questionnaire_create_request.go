package request

import (
	"github.com/QBG-P2/Voting-System/internal/models"
	"time"
)

type QuestionnaireCreateRequest struct {
	models.BaseValidator `json:"-"`
	Title                string                    `json:"title" validate:"required"`
	Description          string                    `json:"description"`
	StartTime            time.Time                 `json:"start_time" validate:"required"`
	EndTime              time.Time                 `json:"end_time" validate:"required,gtfield=StartTime"`
	StepType             models.StepType           `json:"step_type" validate:"required,oneof=Sequential Random"`
	AllowBacktrack       bool                      `json:"allow_backtrack"`
	MaxAttempts          int                       `json:"max_attempts" validate:"gte=1"`
	AnonymityLevel       models.TypeAnonymityLevel `json:"anonymity_level" validate:"required,oneof=Public OwnerOnly Anonymous"`
	ResponseEditDeadline time.Time                 `json:"response_edit_deadline"`
	OwnerID              uint                      `json:"owner_id"`
	MinAge               int                       `json:"min_age"`
	MaxAge               int                       `json:"max_age"`
	AllowedGenders       string                    `json:"allowed_genders" validate:"required,oneof=male female all"`
}

func (q *QuestionnaireCreateRequest) Validate() error {
	return q.BaseValidator.Validate(q)
}

func (q *QuestionnaireCreateRequest) MapToModel(model interface{}) error {
	return q.BaseValidator.MapToModel(q, model)
}
