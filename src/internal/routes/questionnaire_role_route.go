package routes

import (
	"github.com/QBG-P2/Voting-System/internal/handlers"
	"github.com/QBG-P2/Voting-System/internal/repositories"
	"github.com/QBG-P2/Voting-System/internal/services"
)

func questionnaireRoleRoute(rs *RouterService) {

	db := rs.db
	app := rs.app
	cm := rs.questionnaireControlMiddleware

	quesRepo := repositories.NewQuestionnaireRoleRepository(db)
	quesServ := services.NewQuestionnaireRoleService(quesRepo)
	controller := handlers.NewQuestionnaireRoleHandler(quesServ)

	api := app.Group("/api/questionnaire_role")

	api.Get("/questionnaire/:id", cm.CheckPermission("id", "questionnaire-role-GetUserQuestionnaireRoles"), controller.GetUserQuestionnaireRoles)
	api.Post("/assignRole", cm.CheckPermission("questionnaire_id", "questionnaire-role-AssignRole"), controller.AssignRole)

	api.Post("/", cm.CheckPermission("questionnaire_id", "questionnaire-role-CreateQuestionnaireRole"), controller.CreateQuestionnaireRole)
	api.Get("/:id", cm.CheckPermission("questionnaire_id", "questionnaire-role-GetQuestionnaireRole"), controller.GetQuestionnaireRole)
	api.Put("/:id", cm.CheckPermission("questionnaire_id", "questionnaire-role-UpdateQuestionnaireRole"), controller.UpdateQuestionnaireRole)
	api.Delete("/:id", cm.CheckPermission("questionnaire_id", "questionnaire-role-DeleteQuestionnaireRole"), controller.DeleteQuestionnaireRole)

}
