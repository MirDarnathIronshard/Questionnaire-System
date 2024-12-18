package services

import (
	"context"
	"fmt"
	"github.com/QBG-P2/Voting-System/internal/models"
	"github.com/QBG-P2/Voting-System/internal/repositories"
)

type MessageService interface {
	CreateMessage(ctx context.Context, message *models.Message) error
	GetMessagesByChatID(ctx context.Context, chatID uint, offset int, limit int) ([]models.Message, error)
	GetMessageByID(ctx context.Context, id uint) (*models.Message, error)
	UpdateMessage(ctx context.Context, message *models.Message) error
	DeleteMessage(ctx context.Context, id uint) error
	IsOwnerMessage(ctx context.Context, messageID uint, userID uint) (bool, error)
}

func (s *messageService) IsOwnerMessage(ctx context.Context, messageID uint, userID uint) (bool, error) {
	message, err := s.repo.GetByID(ctx, messageID)
	if err != nil {
		return false, err
	}
	if message == nil {
		return false, fmt.Errorf("message not found")
	}
	return message.UserID == userID, nil
}

type messageService struct {
	repo repositories.MessageRepository
}

func NewMessageService(repo repositories.MessageRepository) MessageService {
	return &messageService{
		repo: repo,
	}
}

func (s *messageService) CreateMessage(ctx context.Context, message *models.Message) error {
	if err := message.Validate(); err != nil {
		return err
	}

	return s.repo.Create(ctx, message)
}

func (s *messageService) GetMessagesByChatID(ctx context.Context, chatID uint, offset int, limit int) ([]models.Message, error) {
	return s.repo.GetByChatID(ctx, chatID, offset, limit)
}

func (s *messageService) GetMessageByID(ctx context.Context, id uint) (*models.Message, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *messageService) UpdateMessage(ctx context.Context, message *models.Message) error {
	if err := message.Validate(); err != nil {
		return err
	}

	return s.repo.Update(ctx, message)
}

func (s *messageService) DeleteMessage(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}
