package routes

import (
	"github.com/QBG-P2/Voting-System/internal/handlers"
	"github.com/QBG-P2/Voting-System/internal/repositories"
	"github.com/QBG-P2/Voting-System/internal/services"
)

func questionnaireRolePermissionRoute(rs *RouterService) {
	db := rs.db
	app := rs.app
	cm := rs.questionnaireControlMiddleware

	rolePermRepo := repositories.NewQuestionnaireRolePermissionRepository(db)

	rolePermService := services.NewQuestionnaireRolePermissionService(rolePermRepo)

	handler := handlers.NewQuestionnaireRolePermissionHandler(rolePermService)

	rolePerms := app.Group("/api/questionnaire-role-permissions")

	rolePerms.Post("/assign", cm.CheckPermission("questionnaire_id", "assign-role-permissions"), handler.AssignPermissions)
	rolePerms.Delete("/remove", cm.CheckPermission("questionnaire_id", "remove-role-permissions"), handler.RemovePermissions)
	rolePerms.Get("/:role_id", cm.CheckPermission("questionnaire_id", "view-role-permissions"), handler.GetRolePermissions)
}
