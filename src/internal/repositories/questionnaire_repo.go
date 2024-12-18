package repositories

import (
	"encoding/json"
	"errors"
	"github.com/QBG-P2/Voting-System/internal/models"
	"github.com/QBG-P2/Voting-System/internal/utils"
	"gorm.io/gorm"
	"time"
)

type questionnaireRepository struct {
	db *gorm.DB
}

type QuestionnaireRepository interface {
	Create(questionnaire *models.Questionnaire) error
	GetByID(id uint) (*models.Questionnaire, error)
	Update(questionnaire *models.Questionnaire) error
	Delete(id uint) error
	IsOwner(userID uint, questionnaireID uint) (bool, error)
	UpdateWithOwnership(userID uint, questionnaire *models.Questionnaire) error
	DeleteWithOwnership(userID uint, questionnaireID uint) error
	GetActiveWithOwnership(userID uint) ([]models.Questionnaire, error)
	GetUserQuestionnaires(userID uint) ([]models.Questionnaire, error)
	GetAllActive() ([]models.Questionnaire, error)
	GetAll() ([]models.Questionnaire, error)
	GetPaginated(page int, pageSize int) ([]models.Questionnaire, int, error)
	GetMonitoringData(questionnaireID uint) (string, error)
}

func NewQuestionnaireRepository(db *gorm.DB, cached bool) QuestionnaireRepository {

	return &questionnaireRepository{
		db: db,
	}
}
func (r *questionnaireRepository) Create(questionnaire *models.Questionnaire) error {
	if err := questionnaire.Validate(); err != nil {
		errorsDetail, _ := utils.TranslateError(err, questionnaire)

		return errorsDetail.Error

	}
	return r.db.Create(questionnaire).Error
}

func (r *questionnaireRepository) GetByID(id uint) (*models.Questionnaire, error) {
	var questionnaire models.Questionnaire
	if err := r.db.Preload("Question").
		Preload("QuestionnaireRole").
		Preload("UserQuestionnaireAccess").
		Preload("Chat").
		First(&questionnaire, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("questionnaire not found")
		}
		return nil, err
	}
	return &questionnaire, nil
}

func (r *questionnaireRepository) Update(questionnaire *models.Questionnaire) error {

	if err := questionnaire.Validate(); err != nil {
		errorsDetail, _ := utils.TranslateError(err, questionnaire)

		return errorsDetail.Error

	}
	return r.db.Save(questionnaire).Error
}

func (r *questionnaireRepository) Delete(id uint) error {
	return r.db.Delete(&models.Questionnaire{}, id).Error
}

func (r *questionnaireRepository) IsOwner(userID uint, questionnaireID uint) (bool, error) {
	var questionnaire models.Questionnaire
	if err := r.db.Where("owner_id = ? and id = ?", userID, questionnaireID).First(&questionnaire).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, errors.New("questionnaire not found")
		}
		return false, err
	}
	return questionnaire.OwnerID == userID, nil
}

func (r *questionnaireRepository) UpdateWithOwnership(userID uint, questionnaire *models.Questionnaire) error {
	isOwner, err := r.IsOwner(userID, questionnaire.ID)
	if err != nil {
		return err
	}
	if !isOwner {
		return errors.New("user is not the owner of this questionnaire")
	}

	if err = questionnaire.Validate(); err != nil {
		_, err = utils.TranslateError(err, questionnaire)
		if err != nil {
			return err
		}
	}

	return r.db.Save(questionnaire).Error
}

func (r *questionnaireRepository) DeleteWithOwnership(userID uint, questionnaireID uint) error {
	isOwner, err := r.IsOwner(userID, questionnaireID)
	if err != nil {
		return err
	}
	if !isOwner {
		return errors.New("user is not the owner of this questionnaire")
	}

	return r.db.Delete(&models.Questionnaire{}, questionnaireID).Error
}

func (r *questionnaireRepository) GetActiveWithOwnership(userID uint) ([]models.Questionnaire, error) {
	var list []models.Questionnaire
	currentTime := time.Now()

	err := r.db.Preload("Question").
		Preload("QuestionnaireRole").
		Preload("UserQuestionnaireAccess").
		Preload("Chat").
		Where("owner_id = ? AND end_time > ? AND start_time < ?", userID, currentTime, currentTime).
		Find(&list).Error

	if err != nil {
		return nil, err
	}
	return list, nil
}

func (r *questionnaireRepository) GetUserQuestionnaires(userID uint) ([]models.Questionnaire, error) {
	var list []models.Questionnaire
	err := r.db.Preload("Question").
		Preload("QuestionnaireRole").
		Preload("UserQuestionnaireAccess").
		Preload("Chat").
		Where("owner_id = ?", userID).
		Find(&list).Error

	if err != nil {
		return nil, err
	}
	return list, nil
}

func (r *questionnaireRepository) GetAllActive() ([]models.Questionnaire, error) {
	var list []models.Questionnaire
	currentTime := time.Now()

	err := r.db.Preload("Question").
		Preload("QuestionnaireRole").
		Preload("UserQuestionnaireAccess").
		Preload("Chat").
		Where("end_time > ? AND start_time < ?", currentTime, currentTime).
		Find(&list).Error

	if err != nil {
		return nil, err
	}
	return list, nil
}

func (r *questionnaireRepository) GetAll() ([]models.Questionnaire, error) {
	var list []models.Questionnaire

	err := r.db.Preload("Question").
		Preload("QuestionnaireRole").
		Preload("UserQuestionnaireAccess").
		Preload("Chat").
		Find(&list).Error

	if err != nil {
		return nil, err
	}
	return list, nil
}

func (r *questionnaireRepository) GetPaginated(page int, pageSize int) ([]models.Questionnaire, int, error) {
	var list []models.Questionnaire
	var total int64

	if err := r.db.Model(&models.Questionnaire{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize

	err := r.db.Preload("Question").
		Preload("QuestionnaireRole").
		Preload("UserQuestionnaireAccess").
		Preload("Chat").
		Limit(pageSize).
		Offset(offset).
		Find(&list).Error
	if err != nil {
		return nil, 0, err
	}

	return list, int(total), nil
}

func (r *questionnaireRepository) GetMonitoringData(questionnaireID uint) (string, error) {
	var results []struct {
		UserID    uint      `json:"user_id"`
		Option    string    `json:"option"`
		Timestamp time.Time `json:"timestamp"`
	}

	err := r.db.Raw(`
		SELECT user_id 
		FROM responses
		WHERE questionnaire_id = ?
	`, questionnaireID).Scan(&results).Error

	if err != nil {
		return "", err
	}

	// Format results into a JSON-like string
	formattedData, err := json.Marshal(results)
	if err != nil {
		return "", err
	}

	return string(formattedData), nil
}
