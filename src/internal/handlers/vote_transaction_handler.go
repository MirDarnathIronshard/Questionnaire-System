package handlers

import (
	"github.com/QBG-P2/Voting-System/internal/models"
	"github.com/QBG-P2/Voting-System/internal/services"
	"github.com/QBG-P2/Voting-System/internal/utils"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type CreateTransactionRequest struct {
	models.BaseValidator `json:"-"`
	ToUserID             uint    `json:"to_user_id" validate:"required"`
	Amount               float64 `json:"amount" validate:"required,gt=0"`
}

func NewVoteTransactionHandler(service services.VoteTransactionService) *TransactionHandler {
	return &TransactionHandler{
		transactionService: service,
	}
}

type TransactionHandler struct {
	transactionService services.VoteTransactionService
}

func (h *TransactionHandler) CreateTransaction(c *fiber.Ctx) error {
	var req struct {
		BuyerID         uint    `json:"buyer_id" validate:"required"`
		QuestionnaireID uint    `json:"questionnaire_id" validate:"required"`
		Amount          float64 `json:"amount" validate:"required,gt=0"`
	}
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request", err)
	}

	sellerID, ok := c.Locals("userID").(uint)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid user", nil)
	}

	err := h.transactionService.CreateTransaction(c.UserContext(), sellerID, req.BuyerID, req.QuestionnaireID, req.Amount)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create transaction", err)
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, "Transaction created")
}

func (h *TransactionHandler) ConfirmTransaction(c *fiber.Ctx) error {
	transactionID, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid transaction ID", err)
	}

	err = h.transactionService.ConfirmTransaction(c.UserContext(), uint(transactionID))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to confirm transaction", err)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Transaction confirmed")
}
