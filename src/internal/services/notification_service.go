package services

import (
	"context"
	"encoding/json"
	"github.com/QBG-P2/Voting-System/internal/models"
	"github.com/QBG-P2/Voting-System/internal/repositories"
	"github.com/QBG-P2/Voting-System/pkg/rabbitmq"
	"github.com/QBG-P2/Voting-System/pkg/service_errors"
	"log"
)

type NotificationService interface {
	Create(ctx context.Context, notification *models.Notification) error
	GetByID(ctx context.Context, id uint) (*models.Notification, error)
	GetUserNotifications(ctx context.Context, userID uint, page, pageSize int) ([]models.Notification, int64, error)
	MarkAsRead(ctx context.Context, id uint, userID uint) error
	Delete(ctx context.Context, id uint, userID uint) error
	GetUnreadCount(ctx context.Context, userID uint) (int64, error)
	CreateQuestionnaireNotification(ctx context.Context, userID uint, questionnaireID uint, message string) error
	StartNotificationConsumer(ctx context.Context) error
}

type notificationService struct {
	notificationRepo repositories.NotificationRepository
	rabbitMQ         *rabbitmq.RabbitMQ
}

func NewNotificationService(repo repositories.NotificationRepository, rmq *rabbitmq.RabbitMQ) NotificationService {
	service := &notificationService{
		notificationRepo: repo,
		rabbitMQ:         rmq,
	}

	return service
}

func (s *notificationService) StartNotificationConsumer(ctx context.Context) error {
	return s.rabbitMQ.ConsumeMessages(ctx, "notifications_queue", s.handleNotificationMessage)
}

func (s *notificationService) handleNotificationMessage(msg []byte) error {
	var event struct {
		EventType string `json:"type"`
		Payload   struct {
			Data json.RawMessage `json:"data"`
		} `json:"payload"`
	}

	if err := json.Unmarshal(msg, &event); err != nil {
		return err
	}

	switch event.EventType {
	case "notification.created":
		var notification models.Notification
		if err := json.Unmarshal(event.Payload.Data, &notification); err != nil {
			return err
		}
		//send sms
		log.Printf("Notification created: %v", notification)

	case "notification.updated":
		var notification models.Notification
		if err := json.Unmarshal(event.Payload.Data, &notification); err != nil {
			return err
		}
		log.Printf("Notification updated: %v", notification)

	case "notification.deleted":
		var id uint
		if err := json.Unmarshal(event.Payload.Data, &id); err != nil {
			return err
		}
		log.Printf("Notification deleted: %d", id)
	}

	return nil
}

func (s *notificationService) Create(ctx context.Context, notification *models.Notification) error {
	if notification.UserID == 0 {
		return &service_errors.ServiceError{
			EndUserMessage: "User ID is required for notification",
		}
	}

	if err := notification.Validate(); err != nil {
		return &service_errors.ServiceError{
			EndUserMessage: "Invalid notification data",
			Err:            err,
		}
	}

	return s.notificationRepo.Create(ctx, notification)
}

func (s *notificationService) CreateQuestionnaireNotification(ctx context.Context, userID uint, questionnaireID uint, message string) error {
	notification := &models.Notification{
		UserID:  userID,
		Type:    "questionnaire",
		Message: message,
		IsRead:  false,
	}

	return s.notificationRepo.PublishNotificationEvent(ctx, "notification.questionnaire", notification)
}

func (s *notificationService) GetByID(ctx context.Context, id uint) (*models.Notification, error) {
	notification, err := s.notificationRepo.GetByID(ctx, id)
	if err != nil {
		return nil, &service_errors.ServiceError{
			EndUserMessage: "Notification not found",
			Err:            err,
		}
	}
	return notification, nil
}

func (s *notificationService) GetUserNotifications(ctx context.Context, userID uint, page, pageSize int) ([]models.Notification, int64, error) {
	return s.notificationRepo.GetPaginatedByUserID(ctx, userID, page, pageSize)
}

func (s *notificationService) MarkAsRead(ctx context.Context, id uint, userID uint) error {
	notification, err := s.notificationRepo.GetByID(ctx, id)
	if err != nil {
		return &service_errors.ServiceError{
			EndUserMessage: "Notification not found",
			Err:            err,
		}
	}

	if notification.UserID != userID {
		return &service_errors.ServiceError{
			EndUserMessage: "Unauthorized to modify this notification",
		}
	}

	notification.IsRead = true
	return s.notificationRepo.Update(ctx, notification)
}

func (s *notificationService) Delete(ctx context.Context, id uint, userID uint) error {
	notification, err := s.notificationRepo.GetByID(ctx, id)
	if err != nil {
		return &service_errors.ServiceError{
			EndUserMessage: "Notification not found",
			Err:            err,
		}
	}

	if notification.UserID != userID {
		return &service_errors.ServiceError{
			EndUserMessage: "Unauthorized to delete this notification",
		}
	}

	return s.notificationRepo.Delete(ctx, id)
}

func (s *notificationService) GetUnreadCount(ctx context.Context, userID uint) (int64, error) {
	notifications, err := s.notificationRepo.GetByUserID(ctx, userID)
	if err != nil {
		return 0, err
	}

	var count int64
	for _, notification := range notifications {
		if !notification.IsRead {
			count++
		}
	}
	return count, nil
}
