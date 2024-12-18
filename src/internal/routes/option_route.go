package routes

import (
	"github.com/QBG-P2/Voting-System/internal/handlers"
	"github.com/QBG-P2/Voting-System/internal/repositories"
	"github.com/QBG-P2/Voting-System/internal/services"
)

func optionRoute(rs *RouterService) {
	db := rs.db
	app := rs.app
	cm := rs.questionnaireControlMiddleware

	optionRepo := repositories.NewOptionRepository(db)
	optionService := services.NewOptionService(optionRepo)
	handler := handlers.NewOptionHandler(optionService)
	api := app.Group("/api/options")

	api.Post("/", cm.CheckPermission("questionnaire_id", "create-option"), handler.CreateOption)
	api.Get("/question/:question_id", cm.CheckPermission("questionnaire_id", "get-option-with-question-id"), handler.GetOptionsByQuestionID)
	api.Get("/:id", handler.GetOptionByID)
	api.Put("/:id", cm.CheckPermission("questionnaire_id", "update-option"), handler.UpdateOption)
	api.Delete("/:id", cm.CheckPermission("questionnaire_id", "delete-option"), handler.DeleteOption)
}
