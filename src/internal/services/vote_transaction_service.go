package services

import (
	"context"
	"errors"
	"github.com/QBG-P2/Voting-System/internal/models"
	"github.com/QBG-P2/Voting-System/internal/repositories"
	"github.com/QBG-P2/Voting-System/pkg/auth"
	"gorm.io/gorm"
	"time"
)

type VoteTransactionService struct {
	transactionRepo repositories.VoteTransactionRepository
	questionRepo    repositories.QuestionnaireRepository
	accessRepo      repositories.QuestionnaireAccessRepository
	responseRepo    repositories.ResponseRepository
	userRepo        repositories.UserRepository
}

func NewVoteTransactionService(transactionRepo repositories.VoteTransactionRepository, questionRepo repositories.QuestionnaireRepository, accessRepo repositories.QuestionnaireAccessRepository, responseRepo repositories.ResponseRepository, userRepo repositories.UserRepository) VoteTransactionService {
	return VoteTransactionService{
		transactionRepo: transactionRepo,
		questionRepo:    questionRepo,
		accessRepo:      accessRepo,
		responseRepo:    responseRepo,
		userRepo:        userRepo,
	}
}

func (s *VoteTransactionService) CreateTransaction(ctx context.Context, sellerID, buyerID, questionnaireID uint, amount float64) error {

	questionnaire, err := s.questionRepo.GetByID(questionnaireID)
	if err != nil {
		return err
	}
	_, err = s.accessRepo.GetUserAccess(ctx, sellerID, questionnaireID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("you have not access to this questionnaire")
		}
	}
	if questionnaire.EndTime.Before(time.Now()) {
		return errors.New("questionnaire has already ended")
	}

	transaction := &models.VoteTransaction{
		SellerID:        sellerID,
		BuyerID:         buyerID,
		QuestionnaireID: questionnaireID,
		Amount:          amount,
		Status:          models.TransactionStatusPending,
	}
	err = s.transactionRepo.Create(transaction)
	if err != nil {
		return err
	}

	return nil
}

func (s *VoteTransactionService) ConfirmTransaction(ctx context.Context, transactionID uint) error {
	transaction, err := s.transactionRepo.GetByID(transactionID)
	if err != nil {
		return err
	}

	if transaction.Status != models.TransactionStatusPending {
		return errors.New("invalid transaction status")
	}

	buyer, err := s.userRepo.GetByID(transaction.BuyerID)
	if is, _ := auth.IsSuperAdmin(ctx); !is {
		id, _ := auth.GetUserID(ctx)
		if transaction.BuyerID != *id {
			return errors.New("invalid transaction buyer")
		}
	}
	if err != nil {
		return err
	}
	if buyer.Wallet < transaction.Amount {
		return errors.New("insufficient funds in buyer's wallet")
	}

	err = s.accessRepo.TransferAccess(ctx, transaction.SellerID, transaction.BuyerID, transaction.QuestionnaireID)
	if err != nil {
		return err
	}

	err = s.responseRepo.DeleteByUserAndQuestionnaire(ctx, transaction.SellerID, transaction.QuestionnaireID)
	if err != nil {
		return err
	}

	err = s.userRepo.UpdateWallet(buyer, -transaction.Amount)
	if err != nil {
		return err
	}
	seller, err := s.userRepo.GetByID(transaction.SellerID)
	if err != nil {
		return err
	}
	err = s.userRepo.UpdateWallet(seller, transaction.Amount)
	if err != nil {
		return err
	}

	transaction.Status = models.TransactionStatusConfirmed
	err = s.transactionRepo.Update(transaction)
	if err != nil {
		return err
	}

	return nil
}
