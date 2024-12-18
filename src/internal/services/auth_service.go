package services

import (
	"errors"

	"github.com/casbin/casbin/v2"

	"github.com/QBG-P2/Voting-System/config"
	"github.com/QBG-P2/Voting-System/internal/models"
	"github.com/QBG-P2/Voting-System/internal/repositories"
	"github.com/QBG-P2/Voting-System/internal/utils"
	"github.com/QBG-P2/Voting-System/pkg/security"
)

type AuthService struct {
	userRepo   repositories.UserRepository
	roleRepo   repositories.RoleRepository
	enforcer   *casbin.Enforcer
	JWTSecret  string
	otpService *OtpService
	cfg        *config.Config
}

func NewAuthService(cfg *config.Config, userRepo repositories.UserRepository, roleRepo repositories.RoleRepository, e *casbin.Enforcer, secret string) *AuthService {
	return &AuthService{userRepo: userRepo, roleRepo: roleRepo, JWTSecret: secret, enforcer: e,
		otpService: NewOtpService(cfg), cfg: cfg}
}

func (s *AuthService) RegisterUser(email, password, nationalID, role string, otp string) (*models.User, error) {
	if !utils.ValidateNationalID(nationalID) {
		err := s.otpService.ValidateOtp(email, otp)
		if err != nil {
			return nil, err
		}
		return nil, errors.New("invalid national_id")
	}

	hashedPassword, err := security.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Email:      email,
		Password:   hashedPassword,
		NationalID: nationalID,
		Role:       role,
		Wallet:     100000.0,
	}

	byEmail, _ := s.userRepo.GetUserByEmail(email)
	byNationalID, _ := s.userRepo.GetUserByNationalID(nationalID)

	if byEmail != nil || byNationalID != nil {
		return nil, errors.New("user already exists")
	}

	if err = s.userRepo.CreateUser(user); err != nil {
		return nil, err
	}
	roleName, err := s.roleRepo.GetByName(role)

	if err != nil {
		roleName = &models.Role{
			Name: "user",
		}
		if err = s.roleRepo.Create(roleName); err != nil {
			return nil, err
		}

	}

	if err = s.userRepo.AssignRole(user.ID, roleName.ID); err != nil {
		return nil, err
	}

	if _, err = s.enforcer.AddGroupingPolicy(user.Email, roleName.Name); err != nil {
		return nil, err
	}

	err = s.enforcer.SavePolicy()
	if err != nil {
		return nil, err
	}

	return user, err
}

func (s *AuthService) LoginUser(email, password string) (string, error) {
	user, err := s.userRepo.GetUserByEmail(email)
	if user == nil {
		return "", errors.New("user not found")
	}
	if !security.CheckPasswordHash(password, user.Password) {
		return "", errors.New("invalid credentials")
	}
	roles, err := s.userRepo.GetRoles(user.ID)
	if err != nil {
		return "", err
	}
	return security.GenerateJWTWithUserData(user, s.JWTSecret, roles)
}
