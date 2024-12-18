package repositories

import "github.com/casbin/casbin/v2"

type PolicyRepository interface {
	AddPolicy(role, path, method string) (bool, error)
	RemovePolicy(role, path, method string) (bool, error)
	AddGroupingPolicy(user, role string) (bool, error)
	RemoveGroupingPolicy(user, role string) (bool, error)
	SavePolicy() error
	GetPolicies() ([][]string, error)
}

type policyRepository struct {
	enforcer *casbin.Enforcer
}

func NewPolicyRepository(e *casbin.Enforcer) PolicyRepository {
	return &policyRepository{
		enforcer: e,
	}
}

func (p *policyRepository) AddPolicy(role, path, method string) (bool, error) {
	return p.enforcer.AddPolicy(role, path, method)
}

func (p *policyRepository) RemovePolicy(role, path, method string) (bool, error) {
	return p.enforcer.RemovePolicy(role, path, method)
}

func (p *policyRepository) AddGroupingPolicy(user, role string) (bool, error) {
	return p.enforcer.AddGroupingPolicy(user, role)
}

func (p *policyRepository) RemoveGroupingPolicy(user, role string) (bool, error) {
	return p.enforcer.RemoveGroupingPolicy(user, role)
}

func (p *policyRepository) SavePolicy() error {
	return p.enforcer.SavePolicy()
}

func (p *policyRepository) GetPolicies() ([][]string, error) {
	return p.enforcer.GetPolicy()
}
