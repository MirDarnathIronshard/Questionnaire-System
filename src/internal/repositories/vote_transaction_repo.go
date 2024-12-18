package repositories

import (
	"github.com/QBG-P2/Voting-System/internal/models"
	"gorm.io/gorm"
)

type VoteTransactionRepository interface {
	Create(transaction *models.VoteTransaction) error
	GetByID(id uint) (*models.VoteTransaction, error)
	GetByQuestionID(questionID uint) ([]models.VoteTransaction, error)
	GetByUserID(userID uint) ([]models.VoteTransaction, error)
	GetPaginatedByQuestionID(questionID uint, page int, pageSize int) ([]models.VoteTransaction, int64, error)
	GetPaginatedByUserID(userID uint, page int, pageSize int) ([]models.VoteTransaction, int64, error)
	Update(transaction *models.VoteTransaction) error
	Delete(id uint) error
	DeleteByQuestionID(questionID uint) error
	DeleteByUserID(userID uint) error
}

type voteTransactionRepository struct {
	db *gorm.DB
}

func NewVoteTransactionRepository(db *gorm.DB) VoteTransactionRepository {
	return &voteTransactionRepository{db: db}
}

func (r *voteTransactionRepository) Create(transaction *models.VoteTransaction) error {
	if err := transaction.Validate(); err != nil {
		return err
	}
	return r.db.Create(transaction).Error
}

func (r *voteTransactionRepository) GetByID(id uint) (*models.VoteTransaction, error) {
	var transaction models.VoteTransaction
	if err := r.db.First(&transaction, id).Error; err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (r *voteTransactionRepository) GetByQuestionID(questionID uint) ([]models.VoteTransaction, error) {
	var transactions []models.VoteTransaction
	if err := r.db.Where("question_id = ?", questionID).Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *voteTransactionRepository) GetByUserID(userID uint) ([]models.VoteTransaction, error) {
	var transactions []models.VoteTransaction
	if err := r.db.Where("user_id = ?", userID).Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *voteTransactionRepository) GetPaginatedByQuestionID(questionID uint, page int, pageSize int) ([]models.VoteTransaction, int64, error) {
	var transactions []models.VoteTransaction
	var total int64

	offset := (page - 1) * pageSize

	query := r.db.Where("question_id = ?", questionID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(pageSize).Find(&transactions).Error; err != nil {
		return nil, 0, err
	}

	return transactions, total, nil
}

func (r *voteTransactionRepository) GetPaginatedByUserID(userID uint, page int, pageSize int) ([]models.VoteTransaction, int64, error) {
	var transactions []models.VoteTransaction
	var total int64

	offset := (page - 1) * pageSize

	query := r.db.Where("user_id = ?", userID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(pageSize).Find(&transactions).Error; err != nil {
		return nil, 0, err
	}

	return transactions, total, nil
}

func (r *voteTransactionRepository) Update(transaction *models.VoteTransaction) error {
	if err := transaction.Validate(); err != nil {
		return err
	}
	return r.db.Save(transaction).Error
}

func (r *voteTransactionRepository) Delete(id uint) error {
	return r.db.Delete(&models.VoteTransaction{}, id).Error
}

func (r *voteTransactionRepository) DeleteByQuestionID(questionID uint) error {
	return r.db.Where("question_id = ?", questionID).Delete(&models.VoteTransaction{}).Error
}

func (r *voteTransactionRepository) DeleteByUserID(userID uint) error {
	return r.db.Where("user_id = ?", userID).Delete(&models.VoteTransaction{}).Error
}
