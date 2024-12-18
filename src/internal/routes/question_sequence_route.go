package routes

import (
	"github.com/QBG-P2/Voting-System/internal/handlers"
	"github.com/QBG-P2/Voting-System/internal/repositories"
	"github.com/QBG-P2/Voting-System/internal/services"
	"github.com/QBG-P2/Voting-System/pkg/cache"
)

func questionSequenceRoute(rs *RouterService) {
	db := rs.db
	app := rs.app
	//cm := rs.questionnaireControlMiddleware

	quesRepo := repositories.NewQuestionnaireRepository(db, false)
	questionRepo := repositories.NewQuestionRepository(db)
	resRepo := repositories.NewResponseRepository(db)
	userRepo := repositories.NewUserRepository(db)
	accessRepo := repositories.NewQuestionnaireAccessRepository(db)

	quesServ := services.NewQuestionnaireService(quesRepo, questionRepo, resRepo, rs.notificationService, userRepo, accessRepo)
	sequenceService := services.NewQuestionSequenceService(cache.GetRedis())
	questionService := services.NewQuestionService(questionRepo)

	questionSequence := app.Group("/api/questionnaire-sequence")
	sequenceHandler := handlers.NewQuestionSequenceHandler(questionService, sequenceService, quesServ)

	questionSequence.Post("/:id/start", sequenceHandler.StartQuestionnaire)
	questionSequence.Get("/:id/next", sequenceHandler.GetNextQuestion)
	questionSequence.Get("/:id/current", sequenceHandler.GetCurrentQuestionnaire)
	questionSequence.Get("/:id/previous", sequenceHandler.GetPreviousQuestion)
	questionSequence.Post("/response", sequenceHandler.ValidateAndSubmitResponse)
}
