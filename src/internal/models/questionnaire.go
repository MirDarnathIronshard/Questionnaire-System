package models

import (
	"gorm.io/gorm"
	"time"
)

type TypeAnonymityLevel string
type StepType string

const (
	AnonymityLevelPublic    TypeAnonymityLevel = "Public"
	AnonymityLevelOwnerOnly TypeAnonymityLevel = "OwnerOnly"
	AnonymityLevelAnonymous TypeAnonymityLevel = "Anonymous"
)

const (
	StepTypeSequential StepType = "Sequential"
	StepTypeRandom     StepType = "Random"
)

type QuestionnaireStatus string

const (
	StatusDraft     QuestionnaireStatus = "draft"
	StatusPublished QuestionnaireStatus = "published"
	StatusClosed    QuestionnaireStatus = "closed"
)

type AllowedGenders string

const (
	AllowedGendersMale   AllowedGenders = "male"
	AllowedGendersFemale AllowedGenders = "female"
	AllowedGendersAll    AllowedGenders = "all"
)

type QuestionnaireAnalytics struct {
	QuestionnaireID   uint
	TotalResponses    int
	CompletionRate    float64
	AverageTimeSpent  float64
	ResponsesByOption map[string]int
	DailyResponses    map[string]int
}

type Questionnaire struct {
	gorm.Model
	BaseValidator
	CreatedAt               time.Time                 `json:"created_at"`
	Title                   string                    `gorm:"not null" json:"title" validate:"required"`
	Description             string                    `json:"description"`
	Status                  QuestionnaireStatus       `json:"status"`
	StartTime               time.Time                 `gorm:"not null" json:"start_time" validate:"required"`
	EndTime                 time.Time                 `gorm:"not null" json:"end_time" validate:"required,gtfield=StartTime"`
	StepType                StepType                  `gorm:"not null" json:"step_type" validate:"required,oneof=Sequential Random"`
	AllowBacktrack          bool                      `gorm:"default:true" json:"allow_backtrack"`
	MaxAttempts             int                       `gorm:"default:1" json:"max_attempts"`
	AnonymityLevel          TypeAnonymityLevel        `gorm:"not null" json:"anonymity_level" validate:"required,oneof=Public OwnerOnly Anonymous"`
	ResponseEditDeadline    time.Time                 `json:"response_edit_deadline"`
	OwnerID                 uint                      `gorm:"not null" json:"owner_id" validate:"required"`
	Owner                   User                      `gorm:"foreignKey:OwnerID" json:"-"`
	Question                []Question                `gorm:"foreignKey:QuestionnaireID" json:"questions"`
	QuestionnaireRole       []QuestionnaireRole       `gorm:"foreignKey:QuestionnaireID" json:"questionnaire_roles"`
	UserQuestionnaireAccess []UserQuestionnaireAccess `gorm:"foreignKey:QuestionnaireID" json:"user_questionnaire_access"`
	Chat                    Chat                      `json:"-"`
	MinAge                  int                       `json:"min_age"`
	MaxAge                  int                       `json:"max_age"`
	AllowedGenders          AllowedGenders            `json:"allowed_genders" validate:"required,oneof=male female all"`
}

func (q *Questionnaire) Validate() error {
	return q.BaseValidator.Validate(q)
}
