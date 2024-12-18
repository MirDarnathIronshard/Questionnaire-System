package routes

import (
	"github.com/QBG-P2/Voting-System/internal/repositories"
	"github.com/QBG-P2/Voting-System/internal/services"
)

func RbacInit(rs *RouterService) {

	roleRepo := repositories.NewRoleRepository(rs.db)
	perRepo := repositories.NewPermissionRepository(rs.db)
	qPerRepo := repositories.NewQuestionnairePermissionRepository(rs.db)
	rbac := services.NewRBACInitService(qPerRepo, roleRepo, perRepo, rs.enforcer)
	err := rbac.InitializeRBAC()
	if err != nil {
		return
	}
	err = rbac.InitQuestionnairePermission()
	if err != nil {
		return
	}
}
