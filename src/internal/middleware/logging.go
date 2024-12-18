package middleware

import (
	"context"
	"github.com/QBG-P2/Voting-System/config"
	"github.com/QBG-P2/Voting-System/pkg/logging"
	"github.com/gofiber/fiber/v2"
	"time"
)

type RequestLog struct {
	UserID       uint          `json:"user_id,omitempty"`
	Method       string        `json:"method"`
	Path         string        `json:"path"`
	IP           string        `json:"ip"`
	StatusCode   int           `json:"status_code"`
	ResponseTime time.Duration `json:"response_time"`
	Error        string        `json:"error,omitempty"`
	HandlerLogs  []HandlerLog  `json:"handler_logs,omitempty"`
	RequestBody  interface{}   `json:"request_body,omitempty"`
	ResponseBody interface{}   `json:"response_body,omitempty"`
	StartTime    time.Time     `json:"start_time"`
}

type HandlerLog struct {
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data,omitempty"`
}

type contextKey string

const (
	loggerKey     contextKey = "logger"
	requestLogKey contextKey = "request_log"
)

const (
	categoryHTTP   logging.Category    = "http"
	subCategoryAPI logging.SubCategory = "api"
)

func LoggerMiddleware(cfg *config.Config) fiber.Handler {
	logger := logging.NewLogger(cfg)
	logger.Init()

	return func(c *fiber.Ctx) error {
		startTime := time.Now()

		// Initialize request log
		requestLog := &RequestLog{
			Method:    c.Method(),
			Path:      c.Path(),
			IP:        c.IP(),
			StartTime: startTime,
		}

		if userID, ok := c.Locals("userID").(uint); ok {
			requestLog.UserID = userID
		}

		var reqBody interface{}
		if err := c.BodyParser(&reqBody); err == nil {
			requestLog.RequestBody = reqBody
		}

		ctx := context.WithValue(c.UserContext(), requestLogKey, requestLog)
		ctx = context.WithValue(ctx, loggerKey, logger)
		c.SetUserContext(ctx)

		err := c.Next()

		requestLog.StatusCode = c.Response().StatusCode()
		requestLog.ResponseTime = time.Since(startTime)

		if err != nil {
			requestLog.Error = err.Error()
		}

		logData := make(map[logging.ExtraKey]interface{})
		logData[logging.ExtraKey("user_id")] = requestLog.UserID
		logData[logging.ExtraKey("method")] = requestLog.Method
		logData[logging.ExtraKey("path")] = requestLog.Path
		logData[logging.ExtraKey("status")] = requestLog.StatusCode
		logData[logging.ExtraKey("duration")] = requestLog.ResponseTime.String()
		logData[logging.ExtraKey("ip")] = requestLog.IP

		if len(requestLog.HandlerLogs) > 0 {
			logData[logging.ExtraKey("handler_logs")] = requestLog.HandlerLogs
		}

		if requestLog.Error != "" {
			logData[logging.ExtraKey("error")] = requestLog.Error
		}

		switch {
		case requestLog.StatusCode >= 500:
			logger.Error(categoryHTTP, subCategoryAPI, "request_completed", logData)
		case requestLog.StatusCode >= 400:
			logger.Warn(categoryHTTP, subCategoryAPI, "request_completed", logData)
		default:
			logger.Info(categoryHTTP, subCategoryAPI, "request_completed", logData)
		}

		return err
	}
}

// AddLogEntry adds a log entry to the current request's log collection
func AddLogEntry(c *fiber.Ctx, message string, data map[string]interface{}) {
	ctx := c.UserContext()
	if requestLog, ok := ctx.Value(requestLogKey).(*RequestLog); ok {
		requestLog.HandlerLogs = append(requestLog.HandlerLogs, HandlerLog{
			Message: message,
			Data:    data,
		})
	}
}

func GetLogger(c *fiber.Ctx) logging.Logger {
	if logger, ok := c.UserContext().Value(loggerKey).(logging.Logger); ok {
		return logger
	}
	return logging.NewLogger(config.GetConfig())
}

func GetRequestLog(c *fiber.Ctx) *RequestLog {
	if requestLog, ok := c.UserContext().Value(requestLogKey).(*RequestLog); ok {
		return requestLog
	}
	return nil
}
