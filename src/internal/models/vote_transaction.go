package models

import (
	"gorm.io/gorm"
)

const (
	TransactionStatusPending   = "pending"
	TransactionStatusConfirmed = "confirmed"
	TransactionStatusCancelled = "cancelled"
)

type VoteTransaction struct {
	gorm.Model
	BaseValidator
	SellerID        uint
	BuyerID         uint
	QuestionnaireID uint
	Amount          float64
	Status          string `gorm:"type:varchar(20);default:'pending'"`
}

func (q *VoteTransaction) Validate() error {
	return q.BaseValidator.Validate(q)
}
