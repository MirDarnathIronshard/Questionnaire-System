package middleware

import (
	"fmt"
	"github.com/QBG-P2/Voting-System/internal/interfaces/http/request"
	"github.com/QBG-P2/Voting-System/internal/models"
	"github.com/QBG-P2/Voting-System/internal/repositories"
	"github.com/QBG-P2/Voting-System/internal/utils"
	"github.com/QBG-P2/Voting-System/pkg/auth"
	"github.com/gofiber/fiber/v2"
	"strconv"
	"time"
)

type QuestionnaireControlMiddleware struct {
	questionnaireRepo repositories.QuestionnaireRepository
	questionRepo      repositories.QuestionRepository
	responseRepo      repositories.ResponseRepository
	accessRepo        repositories.QuestionnaireAccessRepository
}

type UserProgress struct {
	CurrentQuestionIndex int       `json:"current_question_index"`
	AnsweredQuestions    []uint    `json:"answered_questions"`
	QuestionOrder        []uint    `json:"question_order"`
	StartTime            time.Time `json:"start_time"`
}

func NewQuestionnaireControlMiddleware(
	questionnaireRepo repositories.QuestionnaireRepository,
	questionRepo repositories.QuestionRepository,
	responseRepo repositories.ResponseRepository,
	accessRepo repositories.QuestionnaireAccessRepository,
) *QuestionnaireControlMiddleware {
	return &QuestionnaireControlMiddleware{
		questionnaireRepo: questionnaireRepo,
		questionRepo:      questionRepo,
		responseRepo:      responseRepo,
		accessRepo:        accessRepo,
	}
}

func (m *QuestionnaireControlMiddleware) CheckPermission(QuestionnaireParamName string, permissionName string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if is, _ := auth.IsSuperAdmin(c.UserContext()); is {
			return c.Next()
		}
		questionnaireID := c.Params(QuestionnaireParamName, "")
		if query := c.Query(QuestionnaireParamName, ""); query != "" {
			questionnaireID = query
		}
		userID := c.Locals("userID").(uint)
		if questionnaireID == "" {
			bodyData := make(map[string]interface{})

			if err := c.BodyParser(&bodyData); err != nil {
				err = c.QueryParser(&bodyData)
			}

			if id, exists := bodyData[QuestionnaireParamName]; exists {
				questionnaireID = fmt.Sprintf("%v", id)
			} else {
				return utils.ErrorResponse(c, fiber.StatusBadRequest, "need "+QuestionnaireParamName, nil)
			}
		}

		qID, _ := strconv.Atoi(questionnaireID)

		questionnaire, err := m.questionnaireRepo.GetByID(uint(qID))

		if err != nil {
			return fiber.NewError(fiber.StatusNotFound, "Questionnaire not found")
		}
		if questionnaire.AnonymityLevel == models.AnonymityLevelPublic {
			return c.Next()
		}

		if questionnaire.OwnerID == userID {
			return c.Next()
		}
		now := time.Now()
		if now.Before(questionnaire.StartTime) {
			return fiber.NewError(fiber.StatusForbidden, "Questionnaire has not started yet")
		}
		if now.After(questionnaire.EndTime) {
			return fiber.NewError(fiber.StatusForbidden, "Questionnaire has ended")
		}

		access, err := m.accessRepo.GetUserAccess(c.Context(), userID, uint(qID))
		if err != nil {
			return fiber.NewError(fiber.StatusForbidden, "No access to questionnaire")
		}
		if !access.IsActive || (access.ExpiresAt != nil && now.After(*access.ExpiresAt)) {
			return fiber.NewError(fiber.StatusForbidden, "Access has expired")
		}

		flagControlHavePermission := false
		for _, permission := range access.Role.Permissions {
			if permission.Name == permissionName {
				flagControlHavePermission = true
			}
		}
		if !flagControlHavePermission {
			return fiber.NewError(fiber.StatusForbidden, "No permission "+permissionName)
		}
		return c.Next()
	}
}

func (m *QuestionnaireControlMiddleware) CheckAnonymityLevel() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if isAdmin, _ := auth.IsSuperAdmin(c.UserContext()); isAdmin {
			return c.Next()
		}
		requestorID := c.Locals("userID").(uint)
		var req request.GetResponsesRequest
		if err := c.QueryParser(&req); err != nil {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid query parameters", nil)
		}

		questionnaire, err := m.questionnaireRepo.GetByID(req.QuestionnaireID)
		if err != nil {
			return fiber.NewError(fiber.StatusNotFound, "Questionnaire not found")
		}

		switch questionnaire.AnonymityLevel {
		case models.AnonymityLevelPublic:
			return c.Next()
		case models.AnonymityLevelOwnerOnly:
			isOwner, err := m.questionnaireRepo.IsOwner(requestorID, req.QuestionnaireID)
			if err != nil || !isOwner {
				return fiber.NewError(fiber.StatusForbidden, "Only owner can view detailed responses")
			}
		case models.AnonymityLevelAnonymous:

			if isUser, _ := auth.HasRole(c.UserContext(), "user"); isUser {
				return fiber.NewError(fiber.StatusForbidden, fmt.Sprintf("this questionnaire have Anonymous type can show only for admin and super_admin"))
			}

		}

		return c.Next()
	}
}
