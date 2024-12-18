package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/QBG-P2/Voting-System/internal/interfaces/http/request"
	"github.com/QBG-P2/Voting-System/internal/models"
	"github.com/QBG-P2/Voting-System/internal/services"
	"github.com/QBG-P2/Voting-System/internal/utils"
	"github.com/QBG-P2/Voting-System/pkg/auth"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"log"
	"net/http"
	"strconv"
	"time"
)

type QuestionnaireHandler struct {
	Service *services.QuestionnaireService
}

func NewQuestionnaireHandler(service *services.QuestionnaireService) *QuestionnaireHandler {
	return &QuestionnaireHandler{Service: service}
}

func (qc *QuestionnaireHandler) CreateQuestionnaire(c *fiber.Ctx) error {
	var req request.QuestionnaireCreateRequest
	var q models.Questionnaire

	if err := c.BodyParser(&req); err != nil {
		errList := req.Validate()
		translatedErr, _ := utils.TranslateError(errList, &req)
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error(), translatedErr.ErrorList)
	}

	if err := req.Validate(); err != nil {
		translatedErr, _ := utils.TranslateError(err, &req)

		return utils.ErrorResponse(c, fiber.StatusBadRequest, translatedErr.Message, translatedErr.ErrorList)
	}

	err := req.MapToModel(&q)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Map data error", nil)
	}

	id, err := auth.GetUserID(c.UserContext())

	if err != nil {
		return utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", err)
	}
	q.OwnerID = *id

	err = qc.Service.CreateQuestionnaire(c.UserContext(), &q)

	if err != nil {
		translatedErr, _ := utils.TranslateError(err, &req)
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create questionnaire", translatedErr.ErrorList)
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, q)
}

func (qc *QuestionnaireHandler) GetQuestionnaire(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id <= 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid questionnaire ID", nil)
	}

	questionnaire, err := qc.Service.GetQuestionnaireByID(uint(id))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Not found error", err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, questionnaire)
}

func (qc *QuestionnaireHandler) UpdateQuestionnaire(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id <= 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid questionnaire ID", nil)
	}

	var req request.QuestionnaireUpdateRequest
	var data models.Questionnaire

	if err := c.BodyParser(&req); err != nil {
		errList := req.Validate()
		translatedErr, _ := utils.TranslateError(errList, &req)
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error(), translatedErr.ErrorList)

	}
	err = req.MapToModel(&data)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Map data error", nil)
	}
	if err = req.Validate(); err != nil {
		translatedErr, _ := utils.TranslateError(err, &req)

		return utils.ErrorResponse(c, fiber.StatusBadRequest, translatedErr.Message, translatedErr.ErrorList)
	}

	data.ID = uint(id)

	userID, _ := c.Locals("userID").(uint)
	data.OwnerID = userID
	if isOwner, err := qc.Service.IsOwnerQuestionnaire(data.ID, userID); !isOwner {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "You do not own this questionnaire", err)
	}

	err = qc.Service.UpdateQuestionnaire(c.UserContext(), &data)

	if err != nil {
		errList, _ := utils.TranslateError(err, &req)

		return utils.ErrorResponse(c, fiber.StatusBadRequest, "bad request", errList.ErrorList)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, data)
}

func (qc *QuestionnaireHandler) DeleteQuestionnaire(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(uint)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid user ID", nil)
	}

	id, err := strconv.Atoi(c.Params("id"))
	if _, err = qc.Service.GetQuestionnaireByID(uint(id)); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid questionnaire ID", nil)
	}

	if isOwner, err := qc.Service.IsOwnerQuestionnaire(uint(id), userID); !isOwner {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "You do not own this questionnaire", err)
	}

	if err := qc.Service.DeleteQuestionnaire(uint(id)); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete questionnaire", nil)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "questionnaire is deleted")
}

func (qc *QuestionnaireHandler) GetUserQuestionnaires(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(uint)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid user ID", nil)
	}

	questionnaires, err := qc.Service.GetUserQuestionnaires(userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve questionnaires", nil)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, questionnaires)
}

func (qc *QuestionnaireHandler) GetActiveQuestionnaires(c *fiber.Ctx) error {
	questionnaires, err := qc.Service.GetActiveQuestionnaires()
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve active questionnaires", nil)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, questionnaires)
}

func (qc *QuestionnaireHandler) PublishQuestionnaire(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id <= 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid questionnaire ID", nil)
	}
	err = qc.Service.PublishQuestionnaire(c.UserContext(), uint(id))

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve active questionnaires", nil)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "PublishQuestionnaire")
}
func (qc *QuestionnaireHandler) CanceledQuestionnaire(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id <= 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid questionnaire ID", nil)
	}
	err = qc.Service.CanceledQuestionnaire(c.UserContext(), uint(id))

	if is, _ := auth.IsSuperAdmin(c.UserContext()); !is {
		idUser, err := auth.GetUserID(c.UserContext())
		if err != nil || idUser == nil {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "you are not owner this questionnaire", err)
		}
		_, err = qc.Service.IsOwnerQuestionnaire(uint(id), *idUser)
		if err != nil {
			return err
		}
	}

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "CanceledQuestionnaire")
}

func (qc *QuestionnaireHandler) GetPaginatedQuestionnaires(c *fiber.Ctx) error {
	pageSize, page := utils.ParsePaginationParams(c).PageSize, utils.ParsePaginationParams(c).Page

	questionnaires, total, err := qc.Service.GetPaginatedQuestionnaires(page, pageSize)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve paginated questionnaires", nil)
	}

	totalPages := (total + pageSize - 1) / pageSize
	pagination := models.Pagination{
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: totalPages,
	}

	return utils.PaginatedResponseWrapper(c, fiber.StatusOK, questionnaires, pagination)
}

func (qc *QuestionnaireHandler) GetAllQuestionnaires(c *fiber.Ctx) error {
	questionnaires, err := qc.Service.GetAllQuestionnaires()
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve all questionnaires", nil)
	}
	return utils.SuccessResponse(c, fiber.StatusOK, questionnaires)
}

func (qc *QuestionnaireHandler) Monitoring(c *websocket.Conn) {

	questionnaireIDStr := c.Params("id", "")
	if questionnaireIDStr == "" {
		err := c.WriteMessage(websocket.TextMessage, []byte("Questionnaire ID is required"))
		if err != nil {
			return
		}
		err = c.Close()
		if err != nil {
			return
		}
		return
	}

	questionnaireID, err := strconv.ParseUint(questionnaireIDStr, 10, 64)
	if err != nil || questionnaireID == 0 {
		err := c.WriteMessage(websocket.TextMessage, []byte("Invalid questionnaire ID"))
		if err != nil {
			return
		}
		err = c.Close()
		if err != nil {
			return
		}
		return
	}

	userIDVal := c.Locals("userID")
	userID, ok := userIDVal.(uint)
	if !ok || userID == 0 {
		err := c.WriteMessage(websocket.TextMessage, []byte("Unauthorized"))
		if err != nil {
			return
		}
		err = c.Close()
		if err != nil {
			return
		}
		return
	}

	allowed, err := qc.Service.CanMonitor(uint(questionnaireID), userID)
	if err != nil {
		log.Println("Error checking permissions:", err)
		err := c.WriteMessage(websocket.TextMessage, []byte("Internal server error"))
		if err != nil {
			return
		}
		err = c.Close()
		if err != nil {
			return
		}
		return
	}
	if !allowed {
		log.Println("User does not have permission to monitor questionnaire", questionnaireID)
		err := c.WriteMessage(websocket.TextMessage, []byte("You do not have permission to monitor this questionnaire"))
		if err != nil {
			return
		}
		err = c.Close()
		if err != nil {
			return
		}
		return
	}

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	log.Printf("WebSocket connection established for questionnaire ID: %d\n", questionnaireID)

	for {
		select {
		case t := <-ticker.C:
			data, err := qc.Service.GetAnalytics(context.Background(), uint(questionnaireID))
			if err != nil {
				log.Println("Error getting live updates:", err)
				continue
			}
			marshal, err := json.Marshal(data)
			if err != nil {
				return
			}
			msg := fmt.Sprintf("{\"timestamp\":\"%s\", \"data\":%s}", t.Format(time.RFC3339), marshal)

			err = c.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				log.Println("Error writing message to client:", err)
				err := c.Close()
				if err != nil {
					return
				}
				return
			}
		}
	}
}
