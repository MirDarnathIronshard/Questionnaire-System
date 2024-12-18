package repositories

import (
	"github.com/QBG-P2/Voting-System/internal/models"
	"gorm.io/gorm"
)

type QuestionnaireRolePermissionRepository interface {
	Create(rolePermission *models.QuestionnaireRolePermission) error

	Delete(rolePermission *models.QuestionnaireRolePermission) error

	GetByRoleID(roleID uint) ([]models.QuestionnaireRolePermission, error)

	GetByPermissionID(permissionID uint) ([]models.QuestionnaireRolePermission, error)

	DeleteByRoleID(roleID uint) error

	DeleteByPermissionID(permissionID uint) error

	CreateBulk(rolePermissions []models.QuestionnaireRolePermission) error
}

type questionnaireRolePermissionRepository struct {
	db *gorm.DB
}

func NewQuestionnaireRolePermissionRepository(db *gorm.DB) QuestionnaireRolePermissionRepository {
	return &questionnaireRolePermissionRepository{db: db}
}

func (r *questionnaireRolePermissionRepository) Create(rolePermission *models.QuestionnaireRolePermission) error {
	return r.db.Create(rolePermission).Error
}

func (r *questionnaireRolePermissionRepository) Delete(rolePermission *models.QuestionnaireRolePermission) error {
	return r.db.Delete(rolePermission).Error
}

func (r *questionnaireRolePermissionRepository) GetByRoleID(roleID uint) ([]models.QuestionnaireRolePermission, error) {
	var rolePermissions []models.QuestionnaireRolePermission
	err := r.db.Where("questionnaire_role_id = ?", roleID).Find(&rolePermissions).Error
	return rolePermissions, err
}

func (r *questionnaireRolePermissionRepository) GetByPermissionID(permissionID uint) ([]models.QuestionnaireRolePermission, error) {
	var rolePermissions []models.QuestionnaireRolePermission
	err := r.db.Where("questionnaire_permission_id = ?", permissionID).Find(&rolePermissions).Error
	return rolePermissions, err
}

func (r *questionnaireRolePermissionRepository) DeleteByRoleID(roleID uint) error {
	return r.db.Where("questionnaire_role_id = ?", roleID).Delete(&models.QuestionnaireRolePermission{}).Error
}

func (r *questionnaireRolePermissionRepository) DeleteByPermissionID(permissionID uint) error {
	return r.db.Where("questionnaire_permission_id = ?", permissionID).Delete(&models.QuestionnaireRolePermission{}).Error
}

func (r *questionnaireRolePermissionRepository) CreateBulk(rolePermissions []models.QuestionnaireRolePermission) error {
	return r.db.Create(&rolePermissions).Error
}
