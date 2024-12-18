package routes

import (
	"github.com/casbin/casbin/v2"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/QBG-P2/Voting-System/config"
	"github.com/QBG-P2/Voting-System/internal/middleware"
	"github.com/QBG-P2/Voting-System/internal/repositories"
	"github.com/QBG-P2/Voting-System/internal/services"
	"github.com/QBG-P2/Voting-System/pkg/rabbitmq"
)

type RouterService struct {
	app                            *fiber.App
	db                             *gorm.DB
	enforcer                       *casbin.Enforcer
	rmq                            *rabbitmq.RabbitMQ
	cfg                            *config.Config
	notificationMiddleware         *middleware.NotificationMiddleware
	questionnaireControlMiddleware *middleware.QuestionnaireControlMiddleware
	notificationService            services.NotificationService
}

func NewRouterService(app *fiber.App, db *gorm.DB, cfg *config.Config, rmq *rabbitmq.RabbitMQ, enforcer *casbin.Enforcer) *RouterService {

	repoNotification := repositories.NewNotificationRepository(db, rmq)
	QuestionnaireRepo := repositories.NewQuestionnaireRepository(db, false)
	QuestionRepo := repositories.NewQuestionRepository(db)
	ResponseRepo := repositories.NewResponseRepository(db)
	AccessRepo := repositories.NewQuestionnaireAccessRepository(db)
	notificationRepo := repositories.NewNotificationRepository(db, rmq)
	ns := services.NewNotificationService(notificationRepo, rmq)

	nm, _ := middleware.NewNotificationMiddleware(rmq, repoNotification)
	qcm := middleware.NewQuestionnaireControlMiddleware(QuestionnaireRepo, QuestionRepo, ResponseRepo, AccessRepo)

	return &RouterService{
		app:                            app,
		db:                             db,
		enforcer:                       enforcer,
		rmq:                            rmq,
		cfg:                            cfg,
		notificationMiddleware:         nm,
		questionnaireControlMiddleware: qcm,
		notificationService:            ns,
	}
}

func InitRoutes(app *fiber.App, db *gorm.DB, enforcer *casbin.Enforcer, rmq *rabbitmq.RabbitMQ, cfg *config.Config) {

	routerService := NewRouterService(app, db, cfg, rmq, enforcer)
	RbacInit(routerService)
	//logger
	app.Use(middleware.LoggerMiddleware(cfg))

	// public route
	authRoute(routerService, cfg)
	swaggerRoute(routerService)

	// auth middleware
	app.Use(middleware.JWTMiddleware(cfg.JWT.Secret))
	app.Use(middleware.CasbinMiddleware(enforcer))

	// private route
	questionnaireRolePermissionRoute(routerService)
	userRoutes(routerService)
	roleRoutes(routerService)
	questionnaireAccessRoute(routerService)
	questionnaireRoleRoute(routerService)
	optionRoute(routerService)
	accessControlRoute(routerService)
	chatRoutes(routerService)
	EmailRoute(routerService)
	messageRoute(routerService)
	NotificationRoute(routerService)
	questionRoute(routerService)
	responseRoute(routerService)
	voteTransactionRoute(routerService)
	questionnaireRoute(routerService)
	questionSequenceRoute(routerService)
}
