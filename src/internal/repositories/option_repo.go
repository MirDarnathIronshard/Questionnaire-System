package repositories

import (
	"github.com/QBG-P2/Voting-System/internal/models"
	"gorm.io/gorm"
)

type OptionRepository interface {
	Create(option *models.Option) error
	GetByID(id uint) (*models.Option, error)
	GetByQuestionID(questionID uint) ([]models.Option, error)
	GetPaginatedByQuestionID(questionID uint, page int, pageSize int) ([]models.Option, int64, error)
	Update(option *models.Option) error
	Delete(id uint) error
	DeleteByQuestionID(questionID uint) error
}

type optionRepository struct {
	db *gorm.DB
}

func NewOptionRepository(db *gorm.DB) OptionRepository {
	return &optionRepository{db: db}
}

func (r *optionRepository) Create(option *models.Option) error {
	if err := option.Validate(); err != nil {
		return err
	}
	return r.db.Create(option).Error
}

func (r *optionRepository) GetByID(id uint) (*models.Option, error) {
	var option models.Option
	if err := r.db.First(&option, id).Error; err != nil {
		return nil, err
	}
	return &option, nil
}

func (r *optionRepository) GetByQuestionID(questionID uint) ([]models.Option, error) {
	var options []models.Option
	if err := r.db.Where("question_id = ?", questionID).Find(&options).Error; err != nil {
		return nil, err
	}
	return options, nil
}

func (r *optionRepository) GetPaginatedByQuestionID(questionID uint, page int, pageSize int) ([]models.Option, int64, error) {
	var options []models.Option
	var total int64

	offset := (page - 1) * pageSize

	query := r.db.Where("question_id = ?", questionID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(pageSize).Find(&options).Error; err != nil {
		return nil, 0, err
	}

	return options, total, nil
}

func (r *optionRepository) Update(option *models.Option) error {
	if err := option.Validate(); err != nil {
		return err
	}
	return r.db.Save(option).Error
}

func (r *optionRepository) Delete(id uint) error {
	return r.db.Delete(&models.Option{}, id).Error
}

func (r *optionRepository) DeleteByQuestionID(questionID uint) error {
	return r.db.Where("question_id = ?", questionID).Delete(&models.Option{}).Error
}
