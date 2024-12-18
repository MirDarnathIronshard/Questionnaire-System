package handlers

import (
	"github.com/QBG-P2/Voting-System/internal/models"
	"github.com/QBG-P2/Voting-System/internal/services"
	"github.com/QBG-P2/Voting-System/internal/utils"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type QuestionSequenceHandler struct {
	sequenceService      *services.QuestionSequenceService
	questionnaireService *services.QuestionnaireService
	questionService      services.QuestionService
}

func NewQuestionSequenceHandler(questionService services.QuestionService, sequenceService *services.QuestionSequenceService, questionnaireService *services.QuestionnaireService) *QuestionSequenceHandler {
	return &QuestionSequenceHandler{
		sequenceService:      sequenceService,
		questionnaireService: questionnaireService,
		questionService:      questionService,
	}
}

func (h *QuestionSequenceHandler) StartQuestionnaire(c *fiber.Ctx) error {
	questionnaireID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid questionnaire ID", nil)
	}

	userID := c.Locals("userID").(uint)

	questionnaire, err := h.questionnaireService.GetQuestionnaireByID(uint(questionnaireID))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Questionnaire not found", nil)
	}

	questions := make([]uint, len(questionnaire.Question))
	for i, q := range questionnaire.Question {
		questions[i] = q.ID
	}

	sequence, err := h.sequenceService.InitializeSequence(
		c.Context(),
		uint(questionnaireID),
		userID,
		questions,
		questionnaire.StepType == "Random",
		questionnaire.AllowBacktrack,
	)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error(), err)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, sequence)
}
func (h *QuestionSequenceHandler) GetCurrentQuestionnaire(c *fiber.Ctx) error {
	questionnaireID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid questionnaire ID", nil)
	}

	userID := c.Locals("userID").(uint)

	_, err = h.questionnaireService.GetQuestionnaireByID(uint(questionnaireID))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Questionnaire not found", nil)
	}

	sequence, err := h.sequenceService.CurrentSequence(
		c.Context(),
		uint(questionnaireID),
		userID,
	)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error(), err)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, sequence)
}

func (h *QuestionSequenceHandler) GetNextQuestion(c *fiber.Ctx) error {
	questionnaireID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid questionnaire ID", nil)
	}

	userID := c.Locals("userID").(uint)

	questionID, err := h.sequenceService.GetNextQuestion(c.Context(), uint(questionnaireID), userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	question, err := h.questionService.GetQuestionByID(c.Context(), questionID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Question not found", nil)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, question)
}

func (h *QuestionSequenceHandler) GetPreviousQuestion(c *fiber.Ctx) error {
	questionnaireID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid questionnaire ID", nil)
	}

	userID := c.Locals("userID").(uint)

	questionID, err := h.sequenceService.GetPreviousQuestion(c.Context(), uint(questionnaireID), userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	question, err := h.questionService.GetQuestionByID(c.Context(), questionID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Question not found", nil)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, question)
}

func (h *QuestionSequenceHandler) ValidateAndSubmitResponse(c *fiber.Ctx) error {
	var req struct {
		QuestionnaireID uint   `json:"questionnaire_id"`
		QuestionID      uint   `json:"question_id"`
		Answer          string `json:"answer"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}

	userID := c.Locals("userID").(uint)

	err := h.sequenceService.ValidateQuestionSequence(
		c.Context(),
		req.QuestionnaireID,
		userID,
		req.QuestionID,
	)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	response := &models.Response{
		QuestionnaireID: req.QuestionnaireID,
		QuestionID:      req.QuestionID,
		UserID:          userID,
		Answer:          req.Answer,
	}

	err = h.questionnaireService.SubmitResponse(c.Context(), response)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to submit response", err)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, response)
}
