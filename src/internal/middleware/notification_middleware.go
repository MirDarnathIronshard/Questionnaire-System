package middleware

import (
	"context"
	"encoding/json"
	"github.com/QBG-P2/Voting-System/internal/models"
	"github.com/QBG-P2/Voting-System/internal/repositories"
	"github.com/QBG-P2/Voting-System/pkg/rabbitmq"
	"github.com/gofiber/fiber/v2"
	"strings"
	"time"
)

type NotificationType string

const (
	QuestionnaireCancelled NotificationType = "QUESTIONNAIRE_CANCELLED"
	VoteCancelled          NotificationType = "VOTE_CANCELLED"
	RoleAssigned           NotificationType = "ROLE_ASSIGNED"
)

type NotificationMiddleware struct {
	rmq *rabbitmq.RabbitMQ
}

type NotificationPayload struct {
	Type      NotificationType       `json:"type"`
	UserID    uint                   `json:"user_id"`
	Message   string                 `json:"message"`
	Data      map[string]interface{} `json:"data"`
	CreatedAt time.Time              `json:"created_at"`
}

func NewNotificationMiddleware(rmq *rabbitmq.RabbitMQ, repository repositories.NotificationRepository) (*NotificationMiddleware, error) {

	err := rmq.ConsumeMessages(context.Background(), "notifications_queue", func(msg []byte) error {
		return handleNotification(context.Background(), msg, repository)
	})

	if err != nil {
		return nil, err
	}

	return &NotificationMiddleware{
		rmq: rmq,
	}, nil
}

func (m *NotificationMiddleware) HandleQuestionnaireEvents() fiber.Handler {
	return func(c *fiber.Ctx) error {

		err := c.Next()

		if c.Response().StatusCode() == fiber.StatusOK {
			path := c.Path()
			method := c.Method()
			switch method {
			case "DELETE":
				if strings.Contains(path, "/api/questionnaire/") {
					userID := c.Locals("userID").(uint)
					questionnaireID := c.Params("id")

					notification := &NotificationPayload{
						Type:    QuestionnaireCancelled,
						UserID:  userID,
						Message: "Questionnaire has been deleted",
						Data: map[string]interface{}{
							"questionnaire_id": questionnaireID,
						},
						CreatedAt: time.Now(),
					}

					err = m.PublishNotification(notification)
					if err != nil {
						// todo log
					}
				}
			}
		}

		return err
	}
}

func (m *NotificationMiddleware) QuestionnaireAccessEvent() fiber.Handler {
	return func(c *fiber.Ctx) error {
		originalCtx := c.UserContext()
		ctx := context.WithValue(originalCtx, "notification_handler", m)
		c.SetUserContext(ctx)

		err := c.Next()
		path := c.Path()
		responseCode := c.Response().StatusCode()
		if responseCode == fiber.StatusOK {
			if c.Method() == "POST" && path == "/api/questionnaire-access/assign" {
				var roleAssignment struct {
					UserID uint `json:"user_id"`
					RoleID uint `json:"role_id"`
				}

				if err := c.BodyParser(&roleAssignment); err != nil {
					return err
				}

				notification := &NotificationPayload{
					Type:    RoleAssigned,
					UserID:  roleAssignment.UserID,
					Message: "New role has been assigned to you",
					Data: map[string]interface{}{
						"role_id": roleAssignment.RoleID,
					},
					CreatedAt: time.Now(),
				}

				err := m.PublishNotification(notification)
				if err != nil {
					//todo Log error
				}
			}
		}

		return err
	}
}

func (m *NotificationMiddleware) PublishNotification(notification *NotificationPayload) error {
	notificationBytes, err := json.Marshal(notification)
	if err != nil {
		return err
	}

	err = m.rmq.PublishMessage(context.Background(), string(notification.Type), notificationBytes)
	if err != nil {
		return err
	}

	return nil
}

func handleNotification(ctx context.Context, msg []byte, repo repositories.NotificationRepository) error {

	return nil
	notification := struct {
		Type    string `json:"type"`
		Payload []byte
	}{}
	if err := json.Unmarshal(msg, &notification); err != nil {
		return nil
	}
	if notification.Type == string(RoleAssigned) && notification.Payload != nil {
		data := struct {
			Type    string `json:"type"`
			UserId  int    `json:"user_id"`
			Message string `json:"message"`
			Data    struct {
				RoleId int `json:"role_id"`
			} `json:"data"`
			CreatedAt time.Time `json:"created_at"`
		}{}
		if err := json.Unmarshal(notification.Payload, &data); err != nil {
			return nil
		}
		repo.Create(ctx, &models.Notification{
			UserID:    uint(data.UserId),
			Type:      data.Type,
			Message:   data.Message,
			CreatedAt: data.CreatedAt,
		})
	}
	if notification.Type == string(QuestionnaireCancelled) && notification.Payload != nil {
		data := struct {
			Type    string `json:"type"`
			UserId  int    `json:"user_id"`
			Message string `json:"message"`
			Data    struct {
				QuestionnaireId string `json:"questionnaire_id"`
			} `json:"data"`
			CreatedAt time.Time `json:"created_at"`
		}{}
		if err := json.Unmarshal(notification.Payload, &data); err != nil {
			return nil
		}
		repo.Create(ctx, &models.Notification{
			UserID:    uint(data.UserId),
			Type:      data.Type,
			Message:   data.Message,
			CreatedAt: data.CreatedAt,
		})
	}

	return nil

}
