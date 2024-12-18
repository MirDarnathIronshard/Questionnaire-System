package repositories

import (
	"context"
	"errors"
	"github.com/QBG-P2/Voting-System/internal/models"
	"gorm.io/gorm"
)

type messageRepository struct {
	db *gorm.DB
}

type MessageRepository interface {
	Create(ctx context.Context, message *models.Message) error
	GetByID(ctx context.Context, id uint) (*models.Message, error)
	GetByChatID(ctx context.Context, chatID uint, offset int, limit int) ([]models.Message, error)
	Update(ctx context.Context, message *models.Message) error
	Delete(ctx context.Context, id uint) error
	GetAllByChatID(ctx context.Context, chatID uint) ([]models.Message, error)
	CountByChatID(ctx context.Context, chatID uint) (int64, error)
}

func NewMessageRepository(db *gorm.DB) MessageRepository {
	return &messageRepository{db: db}
}

func (r *messageRepository) Create(ctx context.Context, message *models.Message) error {
	return r.db.WithContext(ctx).Create(message).Error
}

func (r *messageRepository) GetByID(ctx context.Context, id uint) (*models.Message, error) {
	var message models.Message
	err := r.db.WithContext(ctx).
		Preload("User").
		First(&message, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("message not found")
		}
		return nil, err
	}
	return &message, nil
}

func (r *messageRepository) GetByChatID(ctx context.Context, chatID uint, offset int, limit int) ([]models.Message, error) {
	var messages []models.Message
	err := r.db.WithContext(ctx).
		Where("chat_id = ?", chatID).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Preload("User").
		Find(&messages).Error
	return messages, err
}

func (r *messageRepository) GetAllByChatID(ctx context.Context, chatID uint) ([]models.Message, error) {
	var messages []models.Message
	err := r.db.WithContext(ctx).
		Where("chat_id = ?", chatID).
		Order("created_at DESC").
		Preload("User").
		Find(&messages).Error
	return messages, err
}

func (r *messageRepository) Update(ctx context.Context, message *models.Message) error {
	return r.db.WithContext(ctx).Save(message).Error
}

func (r *messageRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Message{}, id).Error
}

func (r *messageRepository) CountByChatID(ctx context.Context, chatID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Message{}).
		Where("chat_id = ?", chatID).
		Count(&count).Error
	return count, err
}
