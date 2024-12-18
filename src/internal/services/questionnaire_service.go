package services

import (
	"bytes"
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/QBG-P2/Voting-System/internal/models"
	"github.com/QBG-P2/Voting-System/internal/repositories"
	"github.com/xuri/excelize/v2"
	"log"
	"strconv"
	"time"
)

type QuestionnaireService struct {
	questionnaireRepo repositories.QuestionnaireRepository
	questionRepo      repositories.QuestionRepository
	responseRepo      repositories.ResponseRepository
	userRepo          repositories.UserRepository
	accessRepo        repositories.QuestionnaireAccessRepository

	notificationSvc NotificationService
}

func NewQuestionnaireService(
	qRepo repositories.QuestionnaireRepository,
	questionRepo repositories.QuestionRepository,
	responseRepo repositories.ResponseRepository,
	notificationSvc NotificationService,
	userRepo repositories.UserRepository,
	accessRepo repositories.QuestionnaireAccessRepository,
) *QuestionnaireService {
	return &QuestionnaireService{
		questionnaireRepo: qRepo,
		questionRepo:      questionRepo,
		responseRepo:      responseRepo,
		notificationSvc:   notificationSvc,
		userRepo:          userRepo,
		accessRepo:        accessRepo,
	}
}

func (s *QuestionnaireService) GetQuestionnaireByID(id uint) (*models.Questionnaire, error) {
	return s.questionnaireRepo.GetByID(id)
}

func (s *QuestionnaireService) GetAllQuestionnaires() ([]models.Questionnaire, error) {
	return s.questionnaireRepo.GetAll()
}

func (s *QuestionnaireService) IsOwnerQuestionnaire(questionnaireID uint, userID uint) (bool, error) {
	if _, err := s.questionnaireRepo.IsOwner(userID, questionnaireID); err != nil {
		return false, err
	}
	return true, nil
}

func (s *QuestionnaireService) DeleteQuestionnaire(id uint) error {
	if id == 0 {
		return fmt.Errorf("invalid questionnaire ID")
	}
	return s.questionnaireRepo.Delete(id)
}

func (s *QuestionnaireService) GetUserQuestionnaires(userID uint) ([]models.Questionnaire, error) {
	return s.questionnaireRepo.GetUserQuestionnaires(userID)
}

func (s *QuestionnaireService) GetActiveQuestionnaires() ([]models.Questionnaire, error) {
	return s.questionnaireRepo.GetAllActive()
}

func (s *QuestionnaireService) GetPaginatedQuestionnaires(page int, pageSize int) ([]models.Questionnaire, int, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return s.questionnaireRepo.GetPaginated(page, pageSize)
}

func (s *QuestionnaireService) GetLiveUpdates(questionnaireID uint) (string, error) {

	questionnaire, err := s.questionnaireRepo.GetByID(questionnaireID)
	if err != nil {
		return "", err
	}

	if questionnaire.AnonymityLevel == models.AnonymityLevelAnonymous {
		return "Anonymous responses cannot be monitored.", nil
	}

	data, err := s.questionnaireRepo.GetMonitoringData(questionnaireID)
	if err != nil {
		return "", err
	}

	return data, nil
}

func (s *QuestionnaireService) CreateQuestionnaire(ctx context.Context, q *models.Questionnaire) error {
	if err := q.Validate(); err != nil {
		return err
	}

	q.Status = models.StatusDraft

	if err := s.questionnaireRepo.Create(q); err != nil {
		return err
	}

	if q.StartTime.Before(time.Now().Add(24 * time.Hour)) {
		if err := s.scheduleNotifications(ctx, q); err != nil {
			log.Fatal(err)
		}
	}

	return nil
}

func (s *QuestionnaireService) PublishQuestionnaire(ctx context.Context, id uint) error {
	q, err := s.questionnaireRepo.GetByID(id)
	if err != nil {
		return err
	}

	if q.Status != models.StatusDraft {
		return errors.New("questionnaire must be in draft status to publish")
	}

	if err = s.validateForPublish(ctx, q); err != nil {
		return err
	}

	q.Status = models.StatusPublished
	if err = s.questionnaireRepo.Update(q); err != nil {
		return err
	}

	if err = s.notifyParticipants(ctx, q); err != nil {
		log.Fatal(err)
	}

	return err
}
func (s *QuestionnaireService) CanceledQuestionnaire(ctx context.Context, id uint) error {
	q, err := s.questionnaireRepo.GetByID(id)
	if err != nil {
		return err
	}

	//if err = s.validateForPublish(ctx, q); err != nil {
	//	return err
	//}

	q.Status = models.StatusClosed
	if err = s.questionnaireRepo.Update(q); err != nil {
		return err
	}

	if err = s.notifyParticipants(ctx, q); err != nil {
		log.Fatal(err)
	}

	return err
}

func (s *QuestionnaireService) UpdateQuestionnaire(ctx context.Context, q *models.Questionnaire) error {
	existing, err := s.questionnaireRepo.GetByID(q.ID)
	if err != nil {
		return err
	}

	if existing.Status == models.StatusPublished {
		return errors.New("cannot update published questionnaire")
	}

	if err := q.Validate(); err != nil {
		return err
	}

	return s.questionnaireRepo.Update(q)
}

func (s *QuestionnaireService) AddQuestion(ctx context.Context, questionnaireID uint, question *models.Question) error {
	q, err := s.questionnaireRepo.GetByID(questionnaireID)
	if err != nil {
		return err
	}

	if q.Status != models.StatusDraft {
		return errors.New("can only add questions to draft questionnaires")
	}

	question.QuestionnaireID = questionnaireID

	count, err := s.questionRepo.GetTotalByQuestionnaireID(ctx, questionnaireID)
	if err != nil {
		return err
	}
	question.Order = count + 1

	if question.IsConditional {
		if err := s.validateConditionalLogic(ctx, question); err != nil {
			return err
		}
	}

	return s.questionRepo.Create(ctx, question)
}

func (s *QuestionnaireService) ReorderQuestions(ctx context.Context, questionnaireID uint, questionOrders map[uint]int) error {
	questions, err := s.questionRepo.GetByQuestionnaireID(ctx, questionnaireID, 1, 1000)
	if err != nil {
		return err
	}

	for _, q := range questions {
		if newOrder, ok := questionOrders[q.ID]; ok {
			q.Order = newOrder
			if err := s.questionRepo.Update(ctx, &q); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *QuestionnaireService) SubmitResponse(ctx context.Context, response *models.Response) error {
	q, err := s.questionnaireRepo.GetByID(response.QuestionnaireID)
	if err != nil {
		return err
	}

	if q.Status != models.StatusPublished {
		return errors.New("questionnaire is not active")
	}

	if _, err = s.canSubmitResponse(ctx, q, response.UserID); err != nil {
		return err
	}

	if err = s.validateResponse(ctx, response); err != nil {
		return err
	}

	response.CreatedAt = time.Now()

	if err = s.responseRepo.Create(ctx, response); err != nil {
		return err
	}

	if _, err = s.updateAnalytics(ctx, q.ID); err != nil {
		log.Fatal(err)
	}

	return nil
}

func (s *QuestionnaireService) GetResponses(ctx context.Context, questionnaireID uint, page, size int) ([]models.Response, error) {
	q, err := s.questionnaireRepo.GetByID(questionnaireID)
	if err != nil {
		return nil, err
	}

	if q.AnonymityLevel == models.AnonymityLevelAnonymous {
		return nil, errors.New("cannot view responses for anonymous questionnaire")
	}

	return s.responseRepo.GetByQuestionnaireID(ctx, questionnaireID, page, size)
}

func (s *QuestionnaireService) GetAnalytics(ctx context.Context, questionnaireID uint) (*models.QuestionnaireAnalytics, error) {
	_, err := s.questionnaireRepo.GetByID(questionnaireID)
	if err != nil {
		return nil, err
	}

	responses, err := s.responseRepo.GetByQuestionnaireID(ctx, questionnaireID, 1, 1000)
	if err != nil {
		return nil, err
	}

	analytics := &models.QuestionnaireAnalytics{
		ResponsesByOption: make(map[string]int),
		DailyResponses:    make(map[string]int),
		QuestionnaireID:   questionnaireID,
	}

	s.calculateAnalytics(responses, analytics)

	return analytics, nil
}

func (s *QuestionnaireService) ExportResponses(ctx context.Context, questionnaireID uint, format string) ([]byte, error) {
	responses, err := s.responseRepo.GetByQuestionnaireID(ctx, questionnaireID, 1, 1000)
	if err != nil {
		return nil, err
	}

	switch format {
	case "csv":
		return s.exportToCSV(responses)
	case "excel":
		return s.exportToExcel(responses)
	default:
		return nil, errors.New("unsupported export format")
	}
}

func (s *QuestionnaireService) validateForPublish(ctx context.Context, q *models.Questionnaire) error {
	questions, err := s.questionRepo.GetByQuestionnaireID(ctx, q.ID, 1, 1000)
	if err != nil {
		return err
	}

	if len(questions) == 0 {
		return errors.New("questionnaire must have at least one question")
	}

	for _, question := range questions {
		if question.IsConditional {
			if err := s.validateConditionalLogic(ctx, &question); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *QuestionnaireService) validateConditionalLogic(ctx context.Context, question *models.Question) error {
	return question.Validate()
}

func (s *QuestionnaireService) validateResponse(ctx context.Context, response *models.Response) error {
	return response.Validate()
}

func (s *QuestionnaireService) canSubmitResponse(ctx context.Context, q *models.Questionnaire, userID uint) (bool, error) {
	now := time.Now()
	if now.Before(q.StartTime) {
		return false, errors.New("questionnaire start time is earlier than questionnaire end time")
	}
	if now.After(q.EndTime) {
		return false, errors.New("questionnaire end time is earlier than questionnaire start time")
	}

	count, err := s.responseRepo.GetUserResponseCount(ctx, userID, q.ID)
	if err != nil {
		return false, errors.New("cannot get user response count")
	}

	if count >= q.MaxAttempts {
		return false, errors.New("maximum attempts reached")
	}

	if !s.userMeetsCriteria(ctx, q, userID) {
		return false, errors.New("user meets criteria for questionnaire")
	}

	return true, nil
}

func (s *QuestionnaireService) userMeetsCriteria(ctx context.Context, q *models.Questionnaire, userID uint) bool {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return false
	}

	if q.MinAge > 0 || q.MaxAge > 0 {
		userAge := calculateAge(user.BirthDate)
		if q.MinAge > 0 && userAge < q.MinAge {
			return false
		}
		if q.MaxAge > 0 && userAge > q.MaxAge {
			return false
		}
	}

	if q.AllowedGenders != models.AllowedGendersAll {
		if user.Gender != string(q.AllowedGenders) {
			return false
		}
	}

	return true
}

func calculateAge(birthdate time.Time) int {
	now := time.Now()
	age := now.Year() - birthdate.Year()
	if now.YearDay() < birthdate.YearDay() {
		age--
	}
	return age
}

func (s *QuestionnaireService) scheduleNotifications(ctx context.Context, q *models.Questionnaire) error {
	// Get all users who have access to the questionnaire
	users, err := s.questionnaireRepo.GetUserQuestionnaires(q.ID)
	if err != nil {
		return err
	}

	// Schedule notifications for each user
	for _, user := range users {
		message := fmt.Sprintf("Questionnaire '%s' will start soon. Get ready!", q.Title)
		err := s.notificationSvc.CreateQuestionnaireNotification(ctx, user.ID, q.ID, message)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *QuestionnaireService) notifyParticipants(ctx context.Context, q *models.Questionnaire) error {
	users, err := s.questionnaireRepo.GetUserQuestionnaires(q.ID)
	if err != nil {
		return err
	}

	for _, user := range users {
		message := ""
		if q.Status == models.StatusClosed {
			message = fmt.Sprintf("Questionnaire '%s' has closed. Participate now!", q.Title)
		}
		if q.Status == models.StatusPublished {
			message = fmt.Sprintf("Questionnaire '%s' has started. Participate now!", q.Title)
		}

		err := s.notificationSvc.CreateQuestionnaireNotification(ctx, user.ID, q.ID, message)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *QuestionnaireService) updateAnalytics(ctx context.Context, questionnaireID uint) (*models.QuestionnaireAnalytics, error) {
	responses, err := s.responseRepo.GetByQuestionnaireID(ctx, questionnaireID, 1, 1000000)
	if err != nil {
		return nil, err
	}

	if len(responses) == 0 {
		return &models.QuestionnaireAnalytics{
			QuestionnaireID:   questionnaireID,
			TotalResponses:    0,
			CompletionRate:    0,
			AverageTimeSpent:  0,
			ResponsesByOption: nil,
			DailyResponses:    nil,
		}, nil
	}

	totalResponses := len(responses)
	completedCount := 0
	var totalTimeSpent time.Duration
	responsesByOption := make(map[string]int)
	dailyResponses := make(map[string]int)

	for _, resp := range responses {
		if resp.IsFinalized {
			completedCount++
		}

		duration := resp.UpdatedAt.Sub(resp.CreatedAt)
		totalTimeSpent += duration

		dateKey := resp.CreatedAt.Format("2006-01-02")
		dailyResponses[dateKey]++

		optionKey := fmt.Sprintf("option_%d", resp.OptionID)
		responsesByOption[optionKey]++
	}

	completionRate := (float64(completedCount) / float64(totalResponses)) * 100.0
	averageTimeSpent := totalTimeSpent.Seconds() / float64(totalResponses)

	analytics := models.QuestionnaireAnalytics{
		QuestionnaireID:   questionnaireID,
		TotalResponses:    totalResponses,
		CompletionRate:    completionRate,
		AverageTimeSpent:  averageTimeSpent,
		ResponsesByOption: responsesByOption,
		DailyResponses:    dailyResponses,
	}

	return &analytics, nil
}

func (s *QuestionnaireService) calculateAnalytics(responses []models.Response, analytics *models.QuestionnaireAnalytics) {
	analytics.TotalResponses = len(responses)

	completedResponses := 0
	totalTimeSpent := time.Duration(0)

	for _, response := range responses {
		if response.IsFinalized {
			completedResponses++
		}

		responseTime := response.UpdatedAt.Sub(response.CreatedAt)
		totalTimeSpent += responseTime

		date := response.CreatedAt.Format("2006-01-02")
		analytics.DailyResponses[date]++

		optionID := strconv.Itoa(int(response.OptionID))
		analytics.ResponsesByOption[optionID]++
	}

	if analytics.TotalResponses > 0 {
		analytics.CompletionRate = float64(completedResponses) / float64(analytics.TotalResponses) * 100
		analytics.AverageTimeSpent = totalTimeSpent.Seconds() / float64(analytics.TotalResponses)
	}
}
func (s *QuestionnaireService) exportToCSV(responses []models.Response) ([]byte, error) {
	records := [][]string{
		{"ID", "UserID", "QuestionnaireID", "QuestionID", "OptionID", "Answer", "IsFinalized", "CreatedAt", "UpdatedAt"},
	}

	for _, response := range responses {
		record := []string{
			strconv.Itoa(int(response.ID)),
			strconv.Itoa(int(response.UserID)),
			strconv.Itoa(int(response.QuestionnaireID)),
			strconv.Itoa(int(response.QuestionID)),
			strconv.Itoa(int(response.OptionID)),
			response.Answer,
			strconv.FormatBool(response.IsFinalized),
			response.CreatedAt.Format(time.RFC3339),
			response.UpdatedAt.Format(time.RFC3339),
		}
		records = append(records, record)
	}

	b := &bytes.Buffer{}
	w := csv.NewWriter(b)

	err := w.WriteAll(records)
	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func (s *QuestionnaireService) exportToExcel(responses []models.Response) ([]byte, error) {
	f := excelize.NewFile()
	sheet := "Responses"
	err := f.SetSheetName("Sheet1", sheet)
	if err != nil {
		return nil, err
	}

	headers := []string{"ID", "UserID", "QuestionnaireID", "QuestionID", "OptionID", "Answer", "IsFinalized", "CreatedAt", "UpdatedAt"}
	err = f.SetSheetRow(sheet, "A1", &headers)
	if err != nil {
		return nil, err
	}

	for i, response := range responses {
		rowNum := i + 2
		err = f.SetCellValue(sheet, fmt.Sprintf("A%d", rowNum), response.ID)
		if err != nil {
			return nil, err
		}
		err = f.SetCellValue(sheet, fmt.Sprintf("B%d", rowNum), response.UserID)
		if err != nil {
			return nil, err
		}
		err = f.SetCellValue(sheet, fmt.Sprintf("C%d", rowNum), response.QuestionnaireID)
		if err != nil {
			return nil, err
		}
		err = f.SetCellValue(sheet, fmt.Sprintf("D%d", rowNum), response.QuestionID)
		if err != nil {
			return nil, err
		}
		err = f.SetCellValue(sheet, fmt.Sprintf("E%d", rowNum), response.OptionID)
		if err != nil {
			return nil, err
		}
		err = f.SetCellValue(sheet, fmt.Sprintf("F%d", rowNum), response.Answer)
		if err != nil {
			return nil, err
		}
		err = f.SetCellValue(sheet, fmt.Sprintf("G%d", rowNum), response.IsFinalized)
		if err != nil {
			return nil, err
		}
		err = f.SetCellValue(sheet, fmt.Sprintf("H%d", rowNum), response.CreatedAt)
		if err != nil {
			return nil, err
		}
		err = f.SetCellValue(sheet, fmt.Sprintf("I%d", rowNum), response.UpdatedAt)
		if err != nil {
			return nil, err
		}
	}

	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (s *QuestionnaireService) CanMonitor(questionnaireID uint, userID uint) (bool, error) {
	q, err := s.questionnaireRepo.GetByID(questionnaireID)
	if err != nil {
		return false, err
	}
	if q == nil {
		return false, errors.New("questionnaire not found")
	}

	if q.OwnerID == userID {
		return true, nil
	}

	hasPermission, err := s.accessRepo.HasPermission(context.Background(), userID, questionnaireID, "view", "monitoring")
	if err != nil {
		return false, err
	}

	if hasPermission {
		return true, nil
	}

	return false, nil
}
