package routes

import (
	"github.com/QBG-P2/Voting-System/internal/handlers"
	"github.com/QBG-P2/Voting-System/internal/repositories"
	"github.com/QBG-P2/Voting-System/internal/services"
	"github.com/gofiber/contrib/websocket"
)

func questionnaireRoute(rs *RouterService) {
	db := rs.db
	app := rs.app
	HandleQuestionnaireEventsMiddleware := rs.notificationMiddleware.HandleQuestionnaireEvents()
	cm := rs.questionnaireControlMiddleware

	quesRepo := repositories.NewQuestionnaireRepository(db, false)
	questionRepo := repositories.NewQuestionRepository(db)
	resRepo := repositories.NewResponseRepository(db)
	userRepo := repositories.NewUserRepository(db)
	accessRepo := repositories.NewQuestionnaireAccessRepository(db)
	quesServ := services.NewQuestionnaireService(quesRepo, questionRepo, resRepo, rs.notificationService, userRepo, accessRepo)
	controller := handlers.NewQuestionnaireHandler(quesServ)

	api := app.Group("/api/questionnaire", HandleQuestionnaireEventsMiddleware)
	api.Get("/", controller.GetPaginatedQuestionnaires)
	api.Post("/", controller.CreateQuestionnaire)
	api.Get("/user", controller.GetUserQuestionnaires)
	api.Get("/active", controller.GetActiveQuestionnaires)
	api.Get("/:id", cm.CheckPermission("id", "showe-questionnaire"), controller.GetQuestionnaire)
	api.Put("/:id", cm.CheckPermission("id", "update-questionnaire"), controller.UpdateQuestionnaire)
	api.Delete("/:id", cm.CheckPermission("id", "delete-questionnaire"), controller.DeleteQuestionnaire)
	api.Get("/monitoring/:id", cm.CheckPermission("id", "monitoring-questionnaire"), cm.CheckAnonymityLevel(), websocket.New(controller.Monitoring))
	api.Post("/publish-questionnaire/:id", cm.CheckPermission("id", "publish-questionnaire"), controller.PublishQuestionnaire)
	api.Post("/canceled-questionnaire/:id", cm.CheckPermission("id", "canceled-questionnaire"), controller.CanceledQuestionnaire)

}
