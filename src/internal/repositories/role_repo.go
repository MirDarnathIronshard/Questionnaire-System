package repositories

import (
	"github.com/QBG-P2/Voting-System/internal/models"
	"gorm.io/gorm"
)

type RoleRepository interface {
	Create(role *models.Role) error
	GetByID(id uint) (*models.Role, error)
	GetByName(name string) (*models.Role, error)
	Update(role *models.Role) error
	Delete(id uint) error
	AssignPermission(roleID, permissionID uint) error
	RemovePermission(roleID, permissionID uint) error
	GetAll() ([]models.Role, error)
}

type roleRepository struct {
	db *gorm.DB
}

func (r *roleRepository) GetAll() ([]models.Role, error) {
	var roles []models.Role

	err := r.db.Preload("Permissions").Find(&roles).Error
	if err != nil {
		return nil, err
	}

	return roles, nil
}

func (r *roleRepository) Create(role *models.Role) error {
	return r.db.Create(role).Error
}

func (r *roleRepository) GetByID(id uint) (*models.Role, error) {
	var role models.Role
	err := r.db.Preload("Permissions").First(&role, id).Error
	return &role, err
}

func (r *roleRepository) GetByName(name string) (*models.Role, error) {
	var role models.Role
	err := r.db.Preload("Permissions").Where("name = ?", name).First(&role).Error
	return &role, err
}

func (r *roleRepository) Update(role *models.Role) error {
	return r.db.Save(role).Error
}

func (r *roleRepository) Delete(id uint) error {
	return r.db.Delete(&models.Role{}, id).Error
}

func (r *roleRepository) AssignPermission(roleID, permissionID uint) error {
	var role models.Role
	if err := r.db.First(&role, roleID).Error; err != nil {
		return err
	}

	var permission models.Permission
	if err := r.db.First(&permission, permissionID).Error; err != nil {
		return err
	}

	return r.db.Model(&role).Association("Permissions").Append(&permission)
}

func (r *roleRepository) RemovePermission(roleID, permissionID uint) error {
	var role models.Role
	if err := r.db.First(&role, roleID).Error; err != nil {
		return err
	}

	var permission models.Permission
	if err := r.db.First(&permission, permissionID).Error; err != nil {
		return err
	}

	return r.db.Model(&role).Association("Permissions").Delete(&permission)
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{db: db}
}
