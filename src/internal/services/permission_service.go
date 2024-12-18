package services

import (
	"github.com/QBG-P2/Voting-System/internal/models"
	"github.com/QBG-P2/Voting-System/internal/repositories"
)

type PermissionService interface {
	Create(permission *models.Permission) error
	GetByID(id uint) (*models.Permission, error)
	GetByName(name string) (*models.Permission, error)
	Update(permission *models.Permission) error
	Delete(id uint) error
}

type permissionService struct {
	permissionRepo repositories.PermissionRepository
}

func NewPermissionService(p repositories.PermissionRepository) PermissionService {
	return &permissionService{
		permissionRepo: p,
	}
}

func (s *permissionService) Create(permission *models.Permission) error {
	return s.permissionRepo.Create(permission)
}

func (s *permissionService) GetByID(id uint) (*models.Permission, error) {
	return s.permissionRepo.GetByID(id)
}

func (s *permissionService) GetByName(name string) (*models.Permission, error) {
	return s.permissionRepo.GetByName(name)
}

func (s *permissionService) Update(permission *models.Permission) error {
	return s.permissionRepo.Update(permission)
}

func (s *permissionService) Delete(id uint) error {
	return s.permissionRepo.Delete(id)
}
