package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"math/rand"
	"time"
)

type QuestionSequenceService struct {
	redisClient *redis.Client
}

func NewQuestionSequenceService(redisClient *redis.Client) *QuestionSequenceService {
	if redisClient == nil {
		panic(errors.New("redis client is nil"))
	}
	return &QuestionSequenceService{
		redisClient: redisClient,
	}
}

// QuestionSequence stores the sequence info for a user's questionnaire attempt
type QuestionSequence struct {
	QuestionnaireID uint   `json:"questionnaire_id"`
	UserID          uint   `json:"user_id"`
	CurrentStep     int    `json:"current_step"`
	QuestionOrder   []uint `json:"question_order"`
	AllowBacktrack  bool   `json:"allow_backtrack"`
}

func (s *QuestionSequenceService) InitializeSequence(ctx context.Context, questionnaireID, userID uint, questions []uint, isRandom, allowBacktrack bool) (*QuestionSequence, error) {
	key := fmt.Sprintf("sequence:%d:%d", questionnaireID, userID)

	// Check if sequence already exists
	exists, err := s.redisClient.Exists(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	if exists == 1 {
		sequence, err := s.getSequence(ctx, questionnaireID, userID)
		if err != nil {
			return nil, err
		}

		return sequence, nil
	}

	questionOrder := make([]uint, len(questions))
	copy(questionOrder, questions)

	if isRandom {
		rand.Shuffle(len(questionOrder), func(i, j int) {
			questionOrder[i], questionOrder[j] = questionOrder[j], questionOrder[i]
		})
	}

	sequence := &QuestionSequence{
		QuestionnaireID: questionnaireID,
		UserID:          userID,
		CurrentStep:     0,
		QuestionOrder:   questionOrder,
		AllowBacktrack:  allowBacktrack,
	}

	err = s.saveSequence(ctx, key, sequence)
	if err != nil {
		return nil, err
	}

	return sequence, nil
}

func (s *QuestionSequenceService) CurrentSequence(ctx context.Context, questionnaireID, userID uint) (*QuestionSequence, error) {
	key := fmt.Sprintf("sequence:%d:%d", questionnaireID, userID)

	exists, err := s.redisClient.Exists(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	if exists == 1 {
		sequence, err := s.getSequence(ctx, questionnaireID, userID)
		if err != nil {
			return nil, err
		}

		return sequence, nil
	}

	return nil, errors.New("sequence not found")
}

func (s *QuestionSequenceService) GetNextQuestion(ctx context.Context, questionnaireID, userID uint) (uint, error) {
	sequence, err := s.getSequence(ctx, questionnaireID, userID)
	if err != nil {
		return 0, err
	}

	if sequence.CurrentStep >= len(sequence.QuestionOrder) {
		return 0, errors.New("no more questions available")
	}

	questionID := sequence.QuestionOrder[sequence.CurrentStep]
	sequence.CurrentStep++

	key := fmt.Sprintf("sequence:%d:%d", questionnaireID, userID)
	err = s.saveSequence(ctx, key, sequence)
	if err != nil {
		return 0, err
	}

	return questionID, nil
}

func (s *QuestionSequenceService) GetPreviousQuestion(ctx context.Context, questionnaireID, userID uint) (uint, error) {
	sequence, err := s.getSequence(ctx, questionnaireID, userID)
	if err != nil {
		return 0, err
	}

	if !sequence.AllowBacktrack {
		return 0, errors.New("backtracking is not allowed for this questionnaire")
	}

	if sequence.CurrentStep <= 1 {
		return 0, errors.New("no previous question available")
	}

	sequence.CurrentStep--
	questionID := sequence.QuestionOrder[sequence.CurrentStep-1]

	key := fmt.Sprintf("sequence:%d:%d", questionnaireID, userID)
	err = s.saveSequence(ctx, key, sequence)
	if err != nil {
		return 0, err
	}

	return questionID, nil
}

func (s *QuestionSequenceService) ValidateQuestionSequence(ctx context.Context, questionnaireID, userID, questionID uint) error {
	sequence, err := s.getSequence(ctx, questionnaireID, userID)
	if err != nil {
		return err
	}

	if sequence.CurrentStep == 0 {
		return errors.New("sequence not started")
	}

	currentQuestionID := sequence.QuestionOrder[sequence.CurrentStep-1]
	if currentQuestionID != questionID {
		return errors.New("invalid question sequence - this is not the current question")
	}

	return nil
}

func (s *QuestionSequenceService) saveSequence(ctx context.Context, key string, sequence *QuestionSequence) error {
	data, err := json.Marshal(sequence)
	if err != nil {
		return fmt.Errorf("failed to marshal sequence: %w", err)
	}
	return s.redisClient.Set(ctx, key, data, 24*time.Hour).Err()
}

func (s *QuestionSequenceService) getSequence(ctx context.Context, questionnaireID, userID uint) (*QuestionSequence, error) {
	key := fmt.Sprintf("sequence:%d:%d", questionnaireID, userID)
	data, err := s.redisClient.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, errors.New("sequence does not exist")
		}
		return nil, fmt.Errorf("failed to get sequence from Redis: %w", err)
	}

	var sequence QuestionSequence
	err = json.Unmarshal(data, &sequence)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal sequence: %w", err)
	}
	return &sequence, nil
}
