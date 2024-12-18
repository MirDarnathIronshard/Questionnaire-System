package repositories

import (
	"github.com/QBG-P2/Voting-System/internal/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return UserRepository{db: db}
}

func (r *UserRepository) CreateUser(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	err := r.db.Where("email = ?", email).First(user).Error
	if err != nil {
		return nil, err
	}
	return user, err
}
func (r *UserRepository) GetUserByNationalID(nationalID string) (*models.User, error) {
	user := &models.User{}
	err := r.db.Where("national_id = ?", nationalID).First(user).Error
	if err != nil {
		return nil, err
	}
	return user, err
}

func (r *UserRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *UserRepository) Delete(id uint) error {
	return r.db.Delete(&models.User{}, id).Error
}

func (r *UserRepository) AssignRole(userID, roleID uint) error {
	var user models.User
	if err := r.db.First(&user, userID).Error; err != nil {
		return err
	}

	var role models.Role
	if err := r.db.First(&role, roleID).Error; err != nil {
		return err
	}

	return r.db.Model(&user).Association("Roles").Append(&role)
}

func (r *UserRepository) RemoveRole(userID, roleID uint) error {
	var user models.User
	if err := r.db.First(&user, userID).Error; err != nil {
		return err
	}

	var role models.Role
	if err := r.db.First(&role, roleID).Error; err != nil {
		return err
	}

	return r.db.Model(&user).Association("Roles").Delete(&role)
}

func (r *UserRepository) GetByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.Preload("Roles.Permissions").First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetRoles(userID uint) ([]models.Role, error) {
	var user models.User
	err := r.db.Preload("Roles").First(&user, userID).Error
	if err != nil {
		return nil, err
	}

	return user.Roles, nil
}

func (r *UserRepository) UpdateWallet(user *models.User, amount float64) error {
	user.Wallet += amount
	return r.Update(user)
}
