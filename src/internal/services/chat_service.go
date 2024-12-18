package services

import (
	"context"
	"errors"

	"github.com/QBG-P2/Voting-System/internal/models"
	"github.com/QBG-P2/Voting-System/internal/repositories"
)

type ChatService struct {
	chatRepo    repositories.ChatRepository
	messageRepo repositories.MessageRepository
	userRepo    repositories.UserRepository
}

func NewChatService(userRepo repositories.UserRepository, messageRepo repositories.MessageRepository, chatRepo repositories.ChatRepository) ChatService {
	return ChatService{userRepo: userRepo, messageRepo: messageRepo, chatRepo: chatRepo}
}

func (s *ChatService) CreateGroupChat(ctx context.Context, questionnaireID uint, ownerID uint) error {
	chat := &models.Chat{
		QuestionnaireID: questionnaireID,
		SenderUserID:    ownerID,
		ReceiverUserID:  ownerID, // For group chat, sender and receiver are same (owner)
		Type:            "group",
		Status:          "active",
	}
	return s.chatRepo.Create(ctx, chat)
}

func (s *ChatService) CreatePrivateChat(ctx context.Context, senderID uint, receiverID uint) error {
	chat := &models.Chat{
		SenderUserID:   senderID,
		ReceiverUserID: receiverID,
		Type:           "private",
		Status:         "active",
	}
	return s.chatRepo.Create(ctx, chat)
}

func (s *ChatService) SendMessage(ctx context.Context, chatID uint, userID uint, content string, attachmentURL *string) error {
	chat, err := s.chatRepo.GetByID(ctx, chatID)
	if err != nil {
		return err
	}

	if chat.Status != "active" {
		return errors.New("chat is not active")
	}

	message := &models.Message{
		ChatID:        chatID,
		UserID:        userID,
		Content:       content,
		AttachmentURL: attachmentURL,
	}

	return s.messageRepo.Create(ctx, message)
}

func (s *ChatService) GetMessages(ctx context.Context, chatID uint, page int, pageSize int) ([]models.Message, error) {
	offset := (page - 1) * pageSize
	return s.messageRepo.GetByChatID(ctx, chatID, offset, pageSize)
}

func (s *ChatService) GetUserChats(ctx context.Context, userID uint) ([]models.Chat, error) {
	return s.chatRepo.GetChatsByUserID(ctx, userID)
}

func (s *ChatService) DeleteMessage(ctx context.Context, messageID uint, userID uint) error {
	message, err := s.messageRepo.GetByID(ctx, messageID)
	if err != nil {
		return err
	}

	if message.UserID != userID {
		return errors.New("unauthorized to delete this message")
	}

	return s.messageRepo.Delete(ctx, messageID)
}
