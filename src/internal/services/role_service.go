package services

import (
	"errors"
	"fmt"
	"github.com/QBG-P2/Voting-System/internal/models"
	"github.com/QBG-P2/Voting-System/internal/repositories"
	"github.com/casbin/casbin/v2"
)

type RoleService interface {
	// Base role operations
	Create(role *models.Role) error
	GetByID(id uint) (*models.Role, error)
	GetByName(name string) (*models.Role, error)
	GetAll() ([]models.Role, error)
	Update(role *models.Role) error
	Delete(id uint) error

	// Permission management
	AssignPermission(roleID, permissionID uint) error
	RemovePermission(roleID, permissionID uint) error

	// User role management
	AssignToUser(userID, roleID uint) error
	RemoveFromUser(userID, roleID uint) error
	GetUserRoles(userID uint) ([]models.Role, error)
	ValidateUserRole(userID uint, roleName string) (bool, error)
}

type roleService struct {
	roleRepo       repositories.RoleRepository
	userRepo       repositories.UserRepository
	permissionRepo repositories.PermissionRepository
	enforcer       *casbin.Enforcer
}

func (s *roleService) GetAll() ([]models.Role, error) {
	var roles []models.Role

	for _, role := range roles {
		roleWithPermissions, err := s.roleRepo.GetByID(role.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get permissions for role %s: %v", role.Name, err)
		}
		role.Permissions = roleWithPermissions.Permissions
	}

	if s.roleRepo == nil {
		return nil, fmt.Errorf("role repository is not initialized")
	}

	roles, err := s.roleRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get roles: %v", err)
	}

	return roles, nil
}

func (s *roleService) AssignToUser(userID, roleID uint) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return fmt.Errorf("failed to find user: %v", err)
	}

	role, err := s.roleRepo.GetByID(roleID)
	if err != nil {
		return fmt.Errorf("failed to find role: %v", err)
	}

	userRoles, err := s.userRepo.GetRoles(userID)
	if err != nil {
		return fmt.Errorf("failed to get user roles: %v", err)
	}

	for _, existingRole := range userRoles {
		if existingRole.ID == roleID {
			return fmt.Errorf("user already has role: %s", role.Name)
		}
	}

	err = s.userRepo.AssignRole(userID, roleID)
	if err != nil {
		return fmt.Errorf("failed to assign role to user: %v", err)
	}

	_, err = s.enforcer.AddGroupingPolicy(user.Email, role.Name)
	if err != nil {
		return fmt.Errorf("failed to add role policy: %v", err)
	}

	return s.enforcer.SavePolicy()
}

func (s *roleService) RemoveFromUser(userID, roleID uint) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return fmt.Errorf("failed to find user: %v", err)
	}

	role, err := s.roleRepo.GetByID(roleID)
	if err != nil {
		return fmt.Errorf("failed to find role: %v", err)
	}

	userRoles, err := s.userRepo.GetRoles(userID)
	if err != nil {
		return fmt.Errorf("failed to get user roles: %v", err)
	}

	hasRole := false
	for _, existingRole := range userRoles {
		if existingRole.ID == roleID {
			hasRole = true
			break
		}
	}

	if !hasRole {
		return fmt.Errorf("user does not have role: %s", role.Name)
	}

	err = s.userRepo.RemoveRole(userID, roleID)
	if err != nil {
		return fmt.Errorf("failed to remove role from user: %v", err)
	}

	_, err = s.enforcer.RemoveGroupingPolicy(user.Email, role.Name)
	if err != nil {
		return fmt.Errorf("failed to remove role policy: %v", err)
	}

	return s.enforcer.SavePolicy()
}

func (s *roleService) GetUserRoles(userID uint) ([]models.Role, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %v", err)
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	return user.Roles, nil
}

func (s *roleService) ValidateUserRole(userID uint, roleName string) (bool, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return false, fmt.Errorf("failed to find user: %v", err)
	}

	role, err := s.roleRepo.GetByName(roleName)
	if err != nil {
		return false, fmt.Errorf("failed to find role: %v", err)
	}

	hasRole, err := s.enforcer.HasGroupingPolicy(user.Email, role.Name)
	if err != nil {
		return false, fmt.Errorf("failed to check role policy: %v", err)
	}

	return hasRole, nil
}

func NewRoleService(userRepo repositories.UserRepository, r repositories.RoleRepository, p repositories.PermissionRepository, e *casbin.Enforcer) RoleService {
	return &roleService{
		roleRepo:       r,
		permissionRepo: p,
		enforcer:       e,
		userRepo:       userRepo,
	}
}

func (s *roleService) Create(role *models.Role) error {
	err := s.roleRepo.Create(role)
	if err != nil {
		return err
	}
	_, err = s.enforcer.AddRoleForUser(role.Name, role.Name)
	return err
}

func (s *roleService) AssignPermission(roleID, permissionID uint) error {
	err := s.roleRepo.AssignPermission(roleID, permissionID)
	if err != nil {
		return err
	}

	role, err := s.roleRepo.GetByID(roleID)
	if err != nil {
		return err
	}

	permission, err := s.permissionRepo.GetByID(permissionID)
	if err != nil {
		return err
	}

	// اضافه کردن سیاست به Casbin
	_, err = s.enforcer.AddPolicy(role.Name, permission.Path, permission.Method)
	if err != nil {
		return err
	}

	return s.enforcer.SavePolicy()
}

func (s *roleService) RemovePermission(roleID, permissionID uint) error {
	err := s.roleRepo.RemovePermission(roleID, permissionID)
	if err != nil {
		return err
	}

	role, err := s.roleRepo.GetByID(roleID)
	if err != nil {
		return err
	}

	permission, err := s.permissionRepo.GetByID(permissionID)
	if err != nil {
		return err
	}

	_, err = s.enforcer.RemovePolicy(role.Name, permission.Path, permission.Method)
	if err != nil {
		return err
	}

	return s.enforcer.SavePolicy()
}

func (s *roleService) GetByID(id uint) (*models.Role, error) {
	return s.roleRepo.GetByID(id)
}

func (s *roleService) GetByName(name string) (*models.Role, error) {
	return s.roleRepo.GetByName(name)
}

func (s *roleService) Update(role *models.Role) error {
	return s.roleRepo.Update(role)
}

func (s *roleService) Delete(id uint) error {
	return s.roleRepo.Delete(id)
}
