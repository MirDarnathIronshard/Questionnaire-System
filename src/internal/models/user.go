package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	BaseValidator
	NationalID   string    `gorm:"unique;not null" json:"national_id"`
	Gender       string    `gorm:"not null" json:"gender"`
	Email        string    `gorm:"unique;not null" json:"email"`
	Password     string    `gorm:"not null" json:"-"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	BirthDate    time.Time `json:"birth_date"`
	City         string    `json:"city"`
	Role         string    `gorm:"not null;default:user"` // admin, owner, user
	Wallet       float64   `gorm:"default: 5000.00"`
	Is2FAEnabled bool      `gorm:"default:false"`
	Message      []Message `gorm:"foreignKey:UserID"`
	//VoteTransaction         []VoteTransaction         `gorm:"foreignKey:FromUserID"`
	Notification            []Notification            `gorm:"foreignkey:UserID"`
	Response                []Response                `gorm:"foreignKey:UserID"`
	UserQuestionnaireAccess []UserQuestionnaireAccess `gorm:"foreignKey:UserID"`
	Roles                   []Role                    `gorm:"many2many:user_roles;" json:"roles"`
}
type TokenDetail struct {
	AccessToken           string
	AccessTokenExpireTime int64
}

func (q *User) Validate() error {
	return q.BaseValidator.Validate(q)
}
