package routes

import (
	"github.com/QBG-P2/Voting-System/internal/handlers"
	"github.com/QBG-P2/Voting-System/internal/repositories"
	"github.com/QBG-P2/Voting-System/internal/services"
)

func questionnaireAccessRoute(rs *RouterService) {
	db := rs.db
	app := rs.app
	QuestionnaireAccessEventMiddleware := rs.notificationMiddleware.QuestionnaireAccessEvent()
	cm := rs.questionnaireControlMiddleware

	acssRepo := repositories.NewQuestionnaireAccessRepository(db)
	qRepo := repositories.NewQuestionnaireRepository(db, false)
	service := services.NewQuestionnaireAccessService(acssRepo, qRepo)
	handler := handlers.NewQuestionnaireAccessHandler(service)

	access := app.Group("/api/questionnaire-access", QuestionnaireAccessEventMiddleware)

	access.Post("/assign", cm.CheckPermission("questionnaire_id", "questionnaire-access-assign"), handler.AssignRole)
	access.Delete("/revoke/:access_id", cm.CheckPermission("questionnaire_id", "questionnaire-access-revoke"), handler.RevokeAccess)
	access.Get("/:questionnaire_id", cm.CheckPermission("questionnaire_id", "questionnaire-access-GetQuestionnaireUsers"), handler.GetQuestionnaireUsers)
}
