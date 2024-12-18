package routes

import (
	"github.com/QBG-P2/Voting-System/internal/handlers"
	"github.com/QBG-P2/Voting-System/internal/repositories"
	"github.com/QBG-P2/Voting-System/internal/services"
)

func voteTransactionRoute(rs *RouterService) {
	db := rs.db
	app := rs.app
	cm := rs.questionnaireControlMiddleware

	transactionRepo := repositories.NewVoteTransactionRepository(db)
	questionRepo := repositories.NewQuestionnaireRepository(db, false)
	accessRepo := repositories.NewQuestionnaireAccessRepository(db)
	responseRepo := repositories.NewResponseRepository(db)
	userRepo := repositories.NewUserRepository(db)

	transactionService := services.NewVoteTransactionService(transactionRepo, questionRepo, accessRepo, responseRepo, userRepo)
	transactionHandler := handlers.NewVoteTransactionHandler(transactionService)

	app.Post("/api/vote-transactions", cm.CheckPermission("questionnaire_id", "vote-transactions"), transactionHandler.CreateTransaction)
	app.Put("/api/vote-transactions/:id/confirm", transactionHandler.ConfirmTransaction)

}
