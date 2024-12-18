package handlers

import (
	"github.com/QBG-P2/Voting-System/internal/interfaces/http/request"
	"github.com/QBG-P2/Voting-System/internal/models"
	"github.com/QBG-P2/Voting-System/internal/services"
	"github.com/QBG-P2/Voting-System/internal/utils"
	"github.com/QBG-P2/Voting-System/pkg/auth"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"strconv"
)

type QuestionnaireRoleHandler struct {
	Service *services.QuestionnaireRoleService
}

func NewQuestionnaireRoleHandler(service *services.QuestionnaireRoleService) *QuestionnaireRoleHandler {
	return &QuestionnaireRoleHandler{Service: service}
}

func (qc *QuestionnaireRoleHandler) CreateQuestionnaireRole(c *fiber.Ctx) error {
	var req request.QuestionnaireRoleCreateRequest
	var model models.QuestionnaireRole
	id, _ := auth.GetUserID(c.UserContext())
	model.UserID = *id

	if err := c.BodyParser(&req); err != nil {
		errList := req.Validate()

		translatedErr, _ := utils.TranslateError(errList, &req)
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error(), translatedErr.ErrorList)
	}

	if err := req.Validate(); err != nil {
		translatedErr, _ := utils.TranslateError(err, &req)

		return utils.ErrorResponse(c, fiber.StatusBadRequest, translatedErr.Message, translatedErr.ErrorList)
	}

	if err := req.MapToModel(&model); err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Map data error", nil)
	}

	if err := qc.Service.CreateQuestionnaireRole(c.UserContext(), &model); err != nil {
		return err
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, model)
}

func (qc *QuestionnaireRoleHandler) GetQuestionnaireRole(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id <= 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid QuestionnaireRole ID", nil)
	}

	QuestionnaireRole, err := qc.Service.GetQuestionnaireRoleByID(c.UserContext(), uint(id))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Not found error", err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, QuestionnaireRole)
}

func (qc *QuestionnaireRoleHandler) UpdateQuestionnaireRole(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id <= 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid QuestionnaireRole ID", nil)
	}

	var req request.QuestionnaireRoleUpdateRequest
	var model models.QuestionnaireRole

	if err := c.BodyParser(&req); err != nil {
		errList := req.Validate()
		translatedErr, _ := utils.TranslateError(errList, &req)
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error(), translatedErr.ErrorList)

	}
	err = req.MapToModel(&model)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Map data error", nil)
	}

	model.ID = uint(id)

	model.UserID, _ = c.Locals("userID").(uint)

	err = qc.Service.UpdateQuestionnaireRole(c.UserContext(), &model)
	if err != nil {

		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error(), err)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, model)
}

func (qc *QuestionnaireRoleHandler) DeleteQuestionnaireRole(c *fiber.Ctx) error {

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id <= 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid QuestionnaireRole ID", nil)
	}

	if err := qc.Service.DeleteQuestionnaireRole(c.UserContext(), uint(id)); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete QuestionnaireRole", nil)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (qc *QuestionnaireRoleHandler) GetUserQuestionnaireRoles(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(uint)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid user ID", nil)
	}
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id <= 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid Questionnaire ID", nil)
	}

	QuestionnaireRoles, err := qc.Service.GetUserQuestionnaireRoles(c.UserContext(), uint(id), userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve QuestionnaireRoles", nil)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, QuestionnaireRoles)
}

func (qc *QuestionnaireRoleHandler) AssignRole(c *fiber.Ctx) error {
	var req request.AssignRoleRequest
	var role models.QuestionnaireRole

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

	isOwnerOrAdmin, err := qc.Service.IsOwnerOrAdmin(c.UserContext(), req.QuestionnaireID, userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error verifying ownership", nil)
	}
	if !isOwnerOrAdmin {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "You do not have permission to assign roles", nil)
	}

	role.QuestionnaireID = req.QuestionnaireID
	role.UserID = req.UserID
	role.Name = req.Role

	if err := qc.Service.AssignRole(c.UserContext(), &role); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to assign role", err)
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, role)
}
