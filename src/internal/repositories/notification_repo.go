package repositories

import (
	"context"
	"github.com/QBG-P2/Voting-System/internal/models"
	"github.com/QBG-P2/Voting-System/pkg/rabbitmq"
	"gorm.io/gorm"
)

type NotificationRepository interface {
	Create(ctx context.Context, notification *models.Notification) error
	GetByID(ctx context.Context, id uint) (*models.Notification, error)
	GetByUserID(ctx context.Context, userID uint) ([]models.Notification, error)
	GetPaginatedByUserID(ctx context.Context, userID uint, page int, pageSize int) ([]models.Notification, int64, error)
	Update(ctx context.Context, notification *models.Notification) error
	Delete(ctx context.Context, id uint) error
	PublishNotificationEvent(ctx context.Context, eventType string, data interface{}) error
}

type notificationRepository struct {
	db       *gorm.DB
	rabbitMQ *rabbitmq.RabbitMQ
}

func NewNotificationRepository(db *gorm.DB, rabbitMQ *rabbitmq.RabbitMQ) NotificationRepository {
	return &notificationRepository{
		db:       db,
		rabbitMQ: rabbitMQ,
	}
}

func (r *notificationRepository) Create(ctx context.Context, notification *models.Notification) error {
	if err := r.db.WithContext(ctx).Create(notification).Error; err != nil {
		return err
	}
	return nil
}

func (r *notificationRepository) GetByID(ctx context.Context, id uint) (*models.Notification, error) {
	var notification models.Notification
	if err := r.db.WithContext(ctx).First(&notification, id).Error; err != nil {
		return nil, err
	}
	return &notification, nil
}

func (r *notificationRepository) GetByUserID(ctx context.Context, userID uint) ([]models.Notification, error) {
	var notifications []models.Notification
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&notifications).Error; err != nil {
		return nil, err
	}
	return notifications, nil
}

func (r *notificationRepository) GetPaginatedByUserID(ctx context.Context, userID uint, page int, pageSize int) ([]models.Notification, int64, error) {
	var notifications []models.Notification
	var total int64

	offset := (page - 1) * pageSize

	query := r.db.Model(&models.Notification{}).WithContext(ctx).Where("user_id = ?", userID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&notifications).Error; err != nil {
		return nil, 0, err
	}

	return notifications, total, nil
}

func (r *notificationRepository) Update(ctx context.Context, notification *models.Notification) error {
	if err := r.db.WithContext(ctx).Save(notification).Error; err != nil {
		return err
	}
	return nil
}

func (r *notificationRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&models.Notification{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (r *notificationRepository) PublishNotificationEvent(ctx context.Context, eventType string, data interface{}) error {
	event := map[string]interface{}{
		"event_type": eventType,
		"data":       data,
	}

	return r.rabbitMQ.PublishMessage(ctx, eventType, event)
}
