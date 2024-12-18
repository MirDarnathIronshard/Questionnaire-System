package services

import (
	"github.com/QBG-P2/Voting-System/internal/repositories"
)

type PolicyService interface {
	AddPolicy(role, path, method string) (bool, error)
	RemovePolicy(role, path, method string) (bool, error)
	AddGroupingPolicy(user, role string) (bool, error)
	RemoveGroupingPolicy(user, role string) (bool, error)
	SavePolicy() error
	GetPolicies() ([][]string, error)
	InitializePolicies() error
}

type policyService struct {
	policyRepo repositories.PolicyRepository
}

func NewPolicyService(p repositories.PolicyRepository) PolicyService {
	return &policyService{
		policyRepo: p,
	}
}

func (p *policyService) AddPolicy(role, path, method string) (bool, error) {
	return p.policyRepo.AddPolicy(role, path, method)
}

func (p *policyService) RemovePolicy(role, path, method string) (bool, error) {
	return p.policyRepo.RemovePolicy(role, path, method)
}

func (p *policyService) AddGroupingPolicy(user, role string) (bool, error) {
	return p.policyRepo.AddGroupingPolicy(user, role)
}

func (p *policyService) RemoveGroupingPolicy(user, role string) (bool, error) {
	return p.policyRepo.RemoveGroupingPolicy(user, role)
}

func (p *policyService) SavePolicy() error {
	return p.policyRepo.SavePolicy()
}

func (p *policyService) GetPolicies() ([][]string, error) {
	return p.policyRepo.GetPolicies()
}

func (p *policyService) InitializePolicies() error {
	// افزودن سیاست‌ها
	_, err := p.AddPolicy("admin", "/users/*", "GET")
	if err != nil {
		return err
	}
	_, err = p.AddPolicy("admin", "/users", "POST")
	if err != nil {
		return err
	}
	// افزودن سیاست‌های گروه‌بندی
	_, err = p.AddGroupingPolicy("alice", "admin")
	if err != nil {
		return err
	}
	// ذخیره سیاست‌ها
	return p.SavePolicy()
}
