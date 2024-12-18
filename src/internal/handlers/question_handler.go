package handlers

import (
	"strconv"

	"github.com/QBG-P2/Voting-System/internal/interfaces/http/request"
	"github.com/QBG-P2/Voting-System/internal/models"
	"github.com/QBG-P2/Voting-System/internal/services"
	"github.com/QBG-P2/Voting-System/internal/utils"
	"github.com/gofiber/fiber/v2"
)

type QuestionHandler struct {
	Service services.QuestionService
}

func NewQuestionHandler(service services.QuestionService) *QuestionHandler {
	return &QuestionHandler{Service: service}
}

func (h *QuestionHandler) CreateQuestion(c *fiber.Ctx) error {
	var req request.CreateQuestionRequest
	var question models.Question

	if err := c.BodyParser(&req); err != nil {
		errList := req.Validate()
		translatedErr, _ := utils.TranslateError(errList, &req)
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", translatedErr.ErrorList)

	}

	if err := req.Validate(); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Validation failed", err)
	}

	userID, ok := c.Locals("userID").(uint)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid user ID", nil)
	}

	isOwner, err := h.Service.IsOwnerQuestionnaire(c.UserContext(), req.QuestionnaireID, userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error verifying questionnaire ownership", nil)
	}
	if !isOwner {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "You do not own this questionnaire", nil)
	}

	question.Text = req.Text
	question.Type = req.Type
	question.IsConditional = req.IsConditional
	question.QuestionnaireID = req.QuestionnaireID
	question.Order = req.Order
	question.Condition = req.Condition
	question.MediaURL = req.MediaURL
	question.CorrectAnswer = req.CorrectAnswer

	for _, opt := range req.Options {
		option := models.Option{
			Text: opt.Text,
		}
		question.Option = append(question.Option, option)
	}

	if err := h.Service.CreateQuestion(c.UserContext(), &question); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create question", err)
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, question)
}

func (h *QuestionHandler) GetQuestionsByQuestionnaireID(c *fiber.Ctx) error {
	var req request.GetQuestionsRequest

	if err := c.QueryParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid query parameters", nil)
	}

	if err := req.Validate(); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Validation failed", err)
	}

	page := req.Pagination.Page
	size := req.Pagination.PageSize

	questions, err := h.Service.GetQuestionsByQuestionnaireID(c.UserContext(), req.QuestionnaireID, page, size)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve questions", nil)
	}

	total, err := h.Service.GetTotalQuestionsByQuestionnaireID(c.UserContext(), req.QuestionnaireID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve total questions", nil)
	}

	totalPages := (total + size - 1) / size

	return utils.PaginatedResponseWrapper(c, fiber.StatusOK, questions, models.Pagination{
		PageSize:   size,
		Page:       page,
		Total:      total,
		TotalPages: totalPages,
	})
}

func (h *QuestionHandler) GetQuestionByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	questionID, err := strconv.Atoi(idParam)
	if err != nil || questionID <= 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid question ID", nil)
	}

	question, err := h.Service.GetQuestionByID(c.UserContext(), uint(questionID))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Question not found", nil)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, question)
}

func (h *QuestionHandler) UpdateQuestion(c *fiber.Ctx) error {
	idParam := c.Params("id")
	questionID, err := strconv.Atoi(idParam)
	if err != nil || questionID <= 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid question ID", nil)
	}

	var req request.UpdateQuestionRequest

	if err := c.BodyParser(&req); err != nil {
		errList := req.Validate()
		translatedErr, _ := utils.TranslateError(errList, &req)
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", translatedErr.ErrorList)

	}

	if err := req.Validate(); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Validation failed", err)
	}

	userID, ok := c.Locals("userID").(uint)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid user ID", nil)
	}

	existingQuestion, err := h.Service.GetQuestionByID(c.UserContext(), uint(questionID))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Question not found", nil)
	}

	isOwner, err := h.Service.IsOwnerQuestionnaire(c.UserContext(), existingQuestion.QuestionnaireID, userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error verifying questionnaire ownership", nil)
	}
	if !isOwner {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "You do not own this questionnaire", nil)
	}

	existingQuestion.Text = req.Text
	existingQuestion.Type = req.Type
	existingQuestion.IsConditional = req.IsConditional
	existingQuestion.Order = req.Order
	existingQuestion.Condition = req.Condition
	existingQuestion.MediaURL = req.MediaURL
	existingQuestion.CorrectAnswer = req.CorrectAnswer

	existingQuestion.Option = []models.Option{}
	for _, opt := range req.Options {
		option := models.Option{
			Text: opt.Text,
		}
		existingQuestion.Option = append(existingQuestion.Option, option)
	}

	if err := h.Service.UpdateQuestion(c.UserContext(), existingQuestion); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update question", err)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, existingQuestion)
}

func (h *QuestionHandler) DeleteQuestion(c *fiber.Ctx) error {
	idParam := c.Params("id")
	questionID, err := strconv.Atoi(idParam)
	if err != nil || questionID <= 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid question ID", nil)
	}

	userID, ok := c.Locals("userID").(uint)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid user ID", nil)
	}

	question, err := h.Service.GetQuestionByID(c.UserContext(), uint(questionID))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Question not found", nil)
	}

	isOwner, err := h.Service.IsOwnerQuestionnaire(c.UserContext(), question.QuestionnaireID, userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error verifying questionnaire ownership", nil)
	}
	if !isOwner {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "You do not own this questionnaire", nil)
	}

	if err := h.Service.DeleteQuestion(c.UserContext(), uint(questionID)); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete question", nil)
	}

	return c.SendStatus(fiber.StatusNoContent)
}
