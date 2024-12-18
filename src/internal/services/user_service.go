package services

import (
	"github.com/QBG-P2/Voting-System/internal/models"
	"github.com/QBG-P2/Voting-System/internal/repositories"
	"github.com/casbin/casbin/v2"
)

type UserService interface {
	GetByID(id uint) (*models.User, error)
	Update(user *models.User) error
	Create(user *models.User) error
	Delete(id uint) error
	AssignRole(userID, roleID uint) error
	RemoveRole(userID, roleID uint) error
	GetByEmail(email string) (*models.User, error)
}

type userService struct {
	userRepo repositories.UserRepository
	roleRepo repositories.RoleRepository
	enforcer *casbin.Enforcer
}

func (s *userService) GetByEmail(email string) (*models.User, error) {
	return s.userRepo.GetUserByEmail(email)
}
func (s *userService) GetByID(id uint) (*models.User, error) {
	return s.userRepo.GetByID(id)
}
func (s *userService) Update(user *models.User) error {
	return s.userRepo.Update(user)
}

func (s *userService) Create(user *models.User) error {
	return s.userRepo.CreateUser(user)
}

func (s *userService) Delete(id uint) error {
	return s.userRepo.Delete(id)
}

func (s *userService) AssignRole(userID, roleID uint) error {
	err := s.userRepo.AssignRole(userID, roleID)
	if err != nil {
		return err
	}

	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}

	role, err := s.roleRepo.GetByID(roleID)
	if err != nil {
		return err
	}

	_, err = s.enforcer.AddGroupingPolicy(user.Email, role.Name)
	if err != nil {
		return err
	}

	return s.enforcer.SavePolicy()
}

func (s *userService) RemoveRole(userID, roleID uint) error {
	err := s.userRepo.RemoveRole(userID, roleID)
	if err != nil {
		return err
	}

	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}

	role, err := s.roleRepo.GetByID(roleID)
	if err != nil {
		return err
	}

	_, err = s.enforcer.RemoveGroupingPolicy(user.Email, role.Name)
	if err != nil {
		return err
	}

	return s.enforcer.SavePolicy()
}
func NewUserService(u repositories.UserRepository, r repositories.RoleRepository, enforcer *casbin.Enforcer) UserService {
	return &userService{
		userRepo: u,
		roleRepo: r,
		enforcer: enforcer,
	}
}
