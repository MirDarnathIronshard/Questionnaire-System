package repositories

import (
	"context"
	"errors"
	"gorm.io/gorm"

	"github.com/QBG-P2/Voting-System/internal/models"
)

type chatRepository struct {
	db *gorm.DB
}

type ChatRepository interface {
	Create(ctx context.Context, chat *models.Chat) error
	GetByID(ctx context.Context, id uint) (*models.Chat, error)
	GetByQuestionnaireID(ctx context.Context, questionnaireID uint) ([]models.Chat, error)
	Update(ctx context.Context, chat *models.Chat) error
	Delete(ctx context.Context, id uint) error
	GetChatsByUserID(ctx context.Context, userID uint) ([]models.Chat, error)
	GetActiveChatsByUserID(ctx context.Context, userID uint) ([]models.Chat, error)
	GetChatBetweenUsers(ctx context.Context, userID1, userID2 uint) (*models.Chat, error)
	CountUserChats(ctx context.Context, userID uint) (int64, error)
}

func NewChatRepository(db *gorm.DB) ChatRepository {
	return &chatRepository{db: db}
}

func (r *chatRepository) Create(ctx context.Context, chat *models.Chat) error {
	return r.db.WithContext(ctx).Create(chat).Error
}

func (r *chatRepository) GetByID(ctx context.Context, id uint) (*models.Chat, error) {
	var chat models.Chat
	err := r.db.WithContext(ctx).
		Preload("Sender").
		Preload("Receiver").
		First(&chat, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("chat not found")
		}
		return nil, err
	}
	return &chat, nil
}

func (r *chatRepository) GetByQuestionnaireID(ctx context.Context, questionnaireID uint) ([]models.Chat, error) {
	var chats []models.Chat
	err := r.db.WithContext(ctx).
		Where("questionnaire_id = ? AND status = ?", questionnaireID, "active").
		Preload("Sender").
		Preload("Receiver").
		Find(&chats).Error
	return chats, err
}

func (r *chatRepository) Update(ctx context.Context, chat *models.Chat) error {
	return r.db.WithContext(ctx).Save(chat).Error
}

func (r *chatRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Model(&models.Chat{}).
		Where("id = ?", id).
		Update("status", "inactive").Error
}

func (r *chatRepository) GetChatsByUserID(ctx context.Context, userID uint) ([]models.Chat, error) {
	var chats []models.Chat
	err := r.db.WithContext(ctx).
		Where("(sender_user_id = ? OR receiver_user_id = ?)", userID, userID).
		Preload("Sender").
		Preload("Receiver").
		Find(&chats).Error
	return chats, err
}

func (r *chatRepository) GetActiveChatsByUserID(ctx context.Context, userID uint) ([]models.Chat, error) {
	var chats []models.Chat
	err := r.db.WithContext(ctx).
		Where("(sender_user_id = ? OR receiver_user_id = ?) AND status = ?", userID, userID, "active").
		Preload("Sender").
		Preload("Receiver").
		Find(&chats).Error
	return chats, err
}

func (r *chatRepository) GetChatBetweenUsers(ctx context.Context, userID1, userID2 uint) (*models.Chat, error) {
	var chat models.Chat
	err := r.db.WithContext(ctx).
		Where("((sender_user_id = ? AND receiver_user_id = ?) OR (sender_user_id = ? AND receiver_user_id = ?)) AND status = ?",
			userID1, userID2, userID2, userID1, "active").
		Preload("Sender").
		Preload("Receiver").
		First(&chat).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("chat not found")
		}
		return nil, err
	}
	return &chat, nil
}

func (r *chatRepository) CountUserChats(ctx context.Context, userID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Chat{}).
		Where("(sender_user_id = ? OR receiver_user_id = ?) AND status = ?", userID, userID, "active").
		Count(&count).Error
	return count, err
}
