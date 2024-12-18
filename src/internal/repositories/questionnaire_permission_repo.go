package repositories

import (
	"github.com/QBG-P2/Voting-System/internal/models"
	"gorm.io/gorm"
)

type QuestionnairePermissionRepository interface {
	Create(permission *models.QuestionnairePermission) error
	GetByID(id uint) (*models.QuestionnairePermission, error)
	GetByQuestionnaireID(questionnaireID uint) ([]models.QuestionnairePermission, error)
	GetByUserID(userID uint) ([]models.QuestionnairePermission, error)
	Update(permission *models.QuestionnairePermission) error
	Delete(id uint) error
	DeleteByQuestionnaireID(questionnaireID uint) error
}

type questionnairePermissionRepository struct {
	db *gorm.DB
}

func NewQuestionnairePermissionRepository(db *gorm.DB) QuestionnairePermissionRepository {
	return &questionnairePermissionRepository{db: db}
}

func (r *questionnairePermissionRepository) Create(permission *models.QuestionnairePermission) error {
	if err := permission.Validate(); err != nil {
		return err
	}
	return r.db.Create(permission).Error
}

func (r *questionnairePermissionRepository) GetByID(id uint) (*models.QuestionnairePermission, error) {
	var permission models.QuestionnairePermission
	if err := r.db.First(&permission, id).Error; err != nil {
		return nil, err
	}
	return &permission, nil
}

func (r *questionnairePermissionRepository) GetByQuestionnaireID(questionnaireID uint) ([]models.QuestionnairePermission, error) {
	var permissions []models.QuestionnairePermission
	if err := r.db.Where("questionnaire_id = ?", questionnaireID).Find(&permissions).Error; err != nil {
		return nil, err
	}
	return permissions, nil
}

func (r *questionnairePermissionRepository) GetByUserID(userID uint) ([]models.QuestionnairePermission, error) {
	var permissions []models.QuestionnairePermission
	if err := r.db.Where("user_id = ?", userID).Find(&permissions).Error; err != nil {
		return nil, err
	}
	return permissions, nil
}

func (r *questionnairePermissionRepository) Update(permission *models.QuestionnairePermission) error {
	if err := permission.Validate(); err != nil {
		return err
	}
	return r.db.Save(permission).Error
}

func (r *questionnairePermissionRepository) Delete(id uint) error {
	return r.db.Delete(&models.QuestionnairePermission{}, id).Error
}

func (r *questionnairePermissionRepository) DeleteByQuestionnaireID(questionnaireID uint) error {
	return r.db.Where("questionnaire_id = ?", questionnaireID).Delete(&models.QuestionnairePermission{}).Error
}
