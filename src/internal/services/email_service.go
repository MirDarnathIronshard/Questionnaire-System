package services

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/QBG-P2/Voting-System/config"
	"github.com/QBG-P2/Voting-System/pkg/rabbitmq"
	"gopkg.in/gomail.v2"
	"log"
)

type EmailService interface {
	SendEmail(ctx context.Context, to string, subject string, body string) error
	SendQuestionnaireInvitation(ctx context.Context, to string, questionnaireTitle string, inviteLink string) error
	SendWelcomeEmail(ctx context.Context, to string, username string) error
	StartEmailConsumer(ctx context.Context) error
}

type emailService struct {
	cfg      *config.Config
	rabbitMQ *rabbitmq.RabbitMQ
}

type EmailMessage struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
	Type    string `json:"type"`
}

func NewEmailService(cfg *config.Config, rabbitMQ *rabbitmq.RabbitMQ) EmailService {
	return &emailService{
		cfg:      cfg,
		rabbitMQ: rabbitMQ,
	}
}

func (s *emailService) SendEmail(ctx context.Context, to string, subject string, body string) error {
	emailMsg := EmailMessage{
		To:      to,
		Subject: subject,
		Body:    body,
		Type:    "general",
	}

	return s.publishEmailMessage(ctx, emailMsg)
}

func (s *emailService) SendQuestionnaireInvitation(ctx context.Context, to string, questionnaireTitle string, inviteLink string) error {
	subject := fmt.Sprintf("Invitation to Participate: %s", questionnaireTitle)
	body := fmt.Sprintf(`
        <h2>You've been invited to participate in a questionnaire</h2>
        <p>You have been invited to participate in: <strong>%s</strong></p>
        <p>Please click the link below to access the questionnaire:</p>
        <a href="%s">Access Questionnaire</a>
    `, questionnaireTitle, inviteLink)

	emailMsg := EmailMessage{
		To:      to,
		Subject: subject,
		Body:    body,
		Type:    "questionnaire_invitation",
	}

	return s.publishEmailMessage(ctx, emailMsg)
}

func (s *emailService) SendWelcomeEmail(ctx context.Context, to string, username string) error {
	subject := "Welcome to Our Platform!"
	body := fmt.Sprintf(`
        <h2>Welcome %s!</h2>
        <p>Thank you for joining our platform. We're excited to have you here!</p>
        <p>You can now:</p>
        <ul>
            <li>Create questionnaires</li>
            <li>Participate in surveys</li>
            <li>View your results</li>
        </ul>
    `, username)

	emailMsg := EmailMessage{
		To:      to,
		Subject: subject,
		Body:    body,
		Type:    "welcome",
	}

	return s.publishEmailMessage(ctx, emailMsg)
}

func (s *emailService) publishEmailMessage(ctx context.Context, msg EmailMessage) error {
	return s.rabbitMQ.PublishMessage(ctx, "email.send", msg)
}

func (s *emailService) StartEmailConsumer(ctx context.Context) error {
	return s.rabbitMQ.ConsumeMessages(ctx, "emails_queue", func(msg []byte) error {
		var emailMsg EmailMessage
		if err := json.Unmarshal(msg, &emailMsg); err != nil {
			return err
		}

		return s.sendEmailImpl(emailMsg)
	})
}

func (s *emailService) sendEmailImpl(msg EmailMessage) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.cfg.Email.From)
	m.SetHeader("To", msg.To)
	m.SetHeader("Subject", msg.Subject)
	m.SetBody("text/html", msg.Body)

	d := gomail.NewDialer(
		s.cfg.Email.SMTPHost,
		s.cfg.Email.SMTPPort,
		s.cfg.Email.Username,
		s.cfg.Email.Password,
	)

	if err := d.DialAndSend(m); err != nil {
		log.Printf("Failed to send email: %v", err)
		return err
	}

	log.Printf("Email sent successfully to %s", msg.To)
	return nil
}
