package handlers

import (
	"github.com/QBG-P2/Voting-System/internal/models"
	"github.com/QBG-P2/Voting-System/internal/services"
	"github.com/QBG-P2/Voting-System/internal/utils"
	"github.com/gofiber/fiber/v2"
)

type FileUploadHandler struct {
	uploadService *services.FileUploadService
}

func NewFileUploadHandler(uploadService *services.FileUploadService) *FileUploadHandler {
	return &FileUploadHandler{
		uploadService: uploadService,
	}
}

func (h *FileUploadHandler) UploadFile(c *fiber.Ctx) error {
	// Get file from request
	file, err := c.FormFile("file")
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "No file provided", nil)
	}

	// Get file type from request
	fileType := c.FormValue("type", "doc")
	var allowedTypes []string

	// Determine allowed file types based on type parameter
	switch fileType {
	case "image":
		allowedTypes = models.AllowedImageTypes
	case "doc":
		allowedTypes = models.AllowedDocTypes
	case "audio":
		allowedTypes = models.AllowedAudioTypes
	case "video":
		allowedTypes = models.AllowedVideoTypes
	default:
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid file type specified", nil)
	}

	// Upload file
	path, err := h.uploadService.Upload(file, allowedTypes)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to upload file", err)
	}

	// Return success response with file URL
	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"url":  h.uploadService.GetURL(path),
		"path": path,
	})
}

func (h *FileUploadHandler) DeleteFile(c *fiber.Ctx) error {
	filePath := c.Params("path")
	if filePath == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "No file path provided", nil)
	}

	if err := h.uploadService.Delete(filePath); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete file", err)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"message": "File deleted successfully",
	})
}

func fileUploadRoute(app *fiber.App, uploadService *services.FileUploadService) {
	handler := NewFileUploadHandler(uploadService)

	api := app.Group("/api/files")
	api.Post("/upload", handler.UploadFile)
	api.Delete("/:path", handler.DeleteFile)
}
