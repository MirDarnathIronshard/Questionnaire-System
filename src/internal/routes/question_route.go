package routes

import (
	"github.com/QBG-P2/Voting-System/internal/handlers"
	"github.com/QBG-P2/Voting-System/internal/repositories"
	"github.com/QBG-P2/Voting-System/internal/services"
)

func questionRoute(rs *RouterService) {
	db := rs.db
	app := rs.app
	cm := rs.questionnaireControlMiddleware

	questionRepo := repositories.NewQuestionRepository(db)
	questionServ := services.NewQuestionService(questionRepo)
	controller := handlers.NewQuestionHandler(questionServ)

	api := app.Group("/api/questions")

	api.Post("/", cm.CheckPermission("questionnaire_id", "create-Question"), controller.CreateQuestion)
	api.Get("/", cm.CheckPermission("questionnaire_id", "getAll-Question"), controller.GetQuestionsByQuestionnaireID)
	api.Get("/:id", cm.CheckPermission("questionnaire_id", "get-Question"), controller.GetQuestionByID)
	api.Put("/:id", cm.CheckPermission("questionnaire_id", "update-Question"), controller.UpdateQuestion)
	api.Delete("/:id", cm.CheckPermission("questionnaire_id", "delete-Question"), controller.DeleteQuestion)
}
