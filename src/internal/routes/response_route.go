package routes

import (
	"github.com/QBG-P2/Voting-System/internal/handlers"
	"github.com/QBG-P2/Voting-System/internal/repositories"
	"github.com/QBG-P2/Voting-System/internal/services"
)

func responseRoute(rs *RouterService) {
	db := rs.db
	app := rs.app

	responseRepo := repositories.NewResponseRepository(db)
	responseServ := services.NewResponseService(responseRepo)
	controller := handlers.NewResponseHandler(responseServ)
	cm := rs.questionnaireControlMiddleware
	api := app.Group("/api/responses")

	api.Post("/", controller.CreateResponse)
	api.Get("/questionnaire", cm.CheckPermission("questionnaire_id", "responses-GetResponsesByQuestionnaireID"), cm.CheckAnonymityLevel(), controller.GetResponsesByQuestionnaireID)
	api.Get("/:id", cm.CheckPermission("questionnaire_id", "responses-GetResponseByID"), controller.GetResponseByID)
	api.Put("/:id", cm.CheckPermission("questionnaire_id", "responses-UpdateResponse"), controller.UpdateResponse)
	api.Delete("/:id", cm.CheckPermission("questionnaire_id", "responses-DeleteResponse"), controller.DeleteResponse)
}
