package tests

import (
	"github.com/QBG-P2/Voting-System/internal/models"
	"github.com/QBG-P2/Voting-System/internal/repositories"
	"github.com/QBG-P2/Voting-System/pkg/test_mock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRepoQuestionnaire(t *testing.T) {
	db := test_mock.SetupTestDB()
	repo := repositories.NewQuestionnaireRepository(db, false)
	t.Run("TestCreateWithAllFields", func(t *testing.T) {
		startTime := time.Now().Add(24 * time.Hour)
		endTime := startTime.Add(7 * 24 * time.Hour)
		responseEditDeadline := endTime.Add(24 * time.Hour)

		questionnaire := &models.Questionnaire{
			Title:                "Comprehensive Survey",
			Description:          "A detailed survey for testing all fields.",
			StartTime:            startTime,
			EndTime:              endTime,
			StepType:             models.StepTypeSequential,
			AllowBacktrack:       false,
			MaxAttempts:          3,
			AnonymityLevel:       models.AnonymityLevelOwnerOnly,
			ResponseEditDeadline: responseEditDeadline,
			OwnerID:              1,
			Chat:                 models.Chat{},
		}

		err := repo.Create(questionnaire)
		assert.NoError(t, err)
		assert.NotZero(t, questionnaire.ID)
	})

	t.Run("TestUpdateWithAllFields", func(t *testing.T) {
		startTime := time.Now().Add(24 * time.Hour)
		endTime := startTime.Add(7 * 24 * time.Hour)
		responseEditDeadline := endTime.Add(24 * time.Hour)

		questionnaire := &models.Questionnaire{
			Title:                "Initial Survey",
			Description:          "Initial Description",
			StartTime:            startTime,
			EndTime:              endTime,
			StepType:             models.StepTypeSequential,
			AllowBacktrack:       true,
			MaxAttempts:          1,
			AnonymityLevel:       models.AnonymityLevelPublic,
			ResponseEditDeadline: responseEditDeadline,
			OwnerID:              1,
			Chat:                 models.Chat{},
		}

		err := repo.Create(questionnaire)
		assert.NoError(t, err)
		assert.NotZero(t, questionnaire.ID)

		newStartTime := startTime.Add(48 * time.Hour)
		newEndTime := newStartTime.Add(14 * 24 * time.Hour)
		newResponseEditDeadline := newEndTime.Add(48 * time.Hour)

		questionnaire.Title = "Updated Survey"
		questionnaire.Description = "Updated Description"
		questionnaire.StartTime = newStartTime
		questionnaire.EndTime = newEndTime
		questionnaire.StepType = models.StepTypeRandom
		questionnaire.AllowBacktrack = false
		questionnaire.MaxAttempts = 5
		questionnaire.AnonymityLevel = models.AnonymityLevelAnonymous
		questionnaire.ResponseEditDeadline = newResponseEditDeadline
		questionnaire.OwnerID = 2

		err = repo.Update(questionnaire)
		assert.NoError(t, err)

		updatedQuestionnaire, err := repo.GetByID(questionnaire.ID)
		assert.NoError(t, err)
		assert.Equal(t, "Updated Survey", updatedQuestionnaire.Title)
		assert.Equal(t, "Updated Description", updatedQuestionnaire.Description)
		assert.Equal(t, newStartTime.In(time.Local), updatedQuestionnaire.StartTime.In(time.Local))
		assert.Equal(t, newEndTime.In(time.Local), updatedQuestionnaire.EndTime.In(time.Local))
		assert.Equal(t, models.StepTypeRandom, updatedQuestionnaire.StepType)
		assert.False(t, updatedQuestionnaire.AllowBacktrack)
		assert.Equal(t, 5, updatedQuestionnaire.MaxAttempts)
		assert.Equal(t, models.AnonymityLevelAnonymous, updatedQuestionnaire.AnonymityLevel)
		assert.Equal(t, newResponseEditDeadline.In(time.Local), updatedQuestionnaire.ResponseEditDeadline.In(time.Local))
		assert.Equal(t, uint(2), updatedQuestionnaire.OwnerID)
	})

	t.Run("TestCreate", func(t *testing.T) {
		questionnaire := &models.Questionnaire{Title: "Test Questionnaire"}
		err := repo.Create(questionnaire)
		assert.Error(t, err)
	})

	t.Run("TestGetByID_NotFound", func(t *testing.T) {
		result, err := repo.GetByID(999)
		assert.Nil(t, result)
		assert.EqualError(t, err, "questionnaire not found")
	})

	t.Run("TestUpdateQuestionnaire", func(t *testing.T) {
		startTime := time.Now().Add(24 * time.Hour)
		endTime := startTime.Add(7 * 24 * time.Hour)
		responseEditDeadline := endTime.Add(24 * time.Hour)

		questionnaire := &models.Questionnaire{
			Title:                "Initial Survey",
			Description:          "Initial Description",
			StartTime:            startTime,
			EndTime:              endTime,
			StepType:             models.StepTypeSequential,
			AllowBacktrack:       true,
			MaxAttempts:          1,
			AnonymityLevel:       models.AnonymityLevelPublic,
			ResponseEditDeadline: responseEditDeadline,
			OwnerID:              1,
			Chat:                 models.Chat{},
		}

		err := repo.Create(questionnaire)

		assert.Nil(t, err)

		questionnaire.Title = "Updated Title"
		err = repo.Update(questionnaire)
		assert.NoError(t, err)

		updated, _ := repo.GetByID(questionnaire.ID)
		assert.Equal(t, "Updated Title", updated.Title)
	})
	t.Run("TestDeleteQuestionnaire", func(t *testing.T) {
		startTime := time.Now().Add(24 * time.Hour)
		endTime := startTime.Add(7 * 24 * time.Hour)
		responseEditDeadline := endTime.Add(24 * time.Hour)

		questionnaire := &models.Questionnaire{
			Title:                "Initial Survey",
			Description:          "Initial Description",
			StartTime:            startTime,
			EndTime:              endTime,
			StepType:             models.StepTypeSequential,
			AllowBacktrack:       true,
			MaxAttempts:          1,
			AnonymityLevel:       models.AnonymityLevelPublic,
			ResponseEditDeadline: responseEditDeadline,
			OwnerID:              1,
			Chat:                 models.Chat{},
		}

		err := repo.Create(questionnaire)
		assert.Nil(t, err)

		err = repo.Delete(questionnaire.ID)
		assert.NoError(t, err)

		result, err := repo.GetByID(questionnaire.ID)
		assert.Nil(t, result)
		assert.EqualError(t, err, "questionnaire not found")
	})
	t.Run("TestIsOwner", func(t *testing.T) {
		startTime := time.Now().Add(24 * time.Hour)
		endTime := startTime.Add(7 * 24 * time.Hour)
		responseEditDeadline := endTime.Add(24 * time.Hour)

		questionnaire := &models.Questionnaire{
			Title:                "Initial Survey",
			Description:          "Initial Description",
			StartTime:            startTime,
			EndTime:              endTime,
			StepType:             models.StepTypeSequential,
			AllowBacktrack:       true,
			MaxAttempts:          1,
			AnonymityLevel:       models.AnonymityLevelPublic,
			ResponseEditDeadline: responseEditDeadline,
			OwnerID:              1,
			Chat:                 models.Chat{},
		}

		err := repo.Create(questionnaire)
		assert.Nil(t, err)

		isOwner, err := repo.IsOwner(1, questionnaire.ID)
		assert.NoError(t, err)
		assert.True(t, isOwner)

		isOwner, err = repo.IsOwner(2, questionnaire.ID)
		assert.NoError(t, err)
		assert.False(t, isOwner)
	})
	t.Run("TestUpdateWithOwnership", func(t *testing.T) {
		startTime := time.Now().Add(24 * time.Hour)
		endTime := startTime.Add(7 * 24 * time.Hour)
		responseEditDeadline := endTime.Add(24 * time.Hour)

		questionnaire := &models.Questionnaire{
			Title:                "Initial Survey",
			Description:          "Initial Description",
			StartTime:            startTime,
			EndTime:              endTime,
			StepType:             models.StepTypeSequential,
			AllowBacktrack:       true,
			MaxAttempts:          1,
			AnonymityLevel:       models.AnonymityLevelPublic,
			ResponseEditDeadline: responseEditDeadline,
			OwnerID:              1,
			Chat:                 models.Chat{},
		}

		err := repo.Create(questionnaire)
		assert.Nil(t, err)

		questionnaire.Title = "Updated Title"
		err = repo.UpdateWithOwnership(1, questionnaire)
		assert.NoError(t, err)

		updated, _ := repo.GetByID(questionnaire.ID)
		assert.Equal(t, "Updated Title", updated.Title)
	})
	t.Run("TestUpdateWithOwnership_NotOwner", func(t *testing.T) {
		startTime := time.Now().Add(24 * time.Hour)
		endTime := startTime.Add(7 * 24 * time.Hour)
		responseEditDeadline := endTime.Add(24 * time.Hour)

		questionnaire := &models.Questionnaire{
			Title:                "Initial Survey",
			Description:          "Initial Description",
			StartTime:            startTime,
			EndTime:              endTime,
			StepType:             models.StepTypeSequential,
			AllowBacktrack:       true,
			MaxAttempts:          1,
			AnonymityLevel:       models.AnonymityLevelPublic,
			ResponseEditDeadline: responseEditDeadline,
			OwnerID:              1,
			Chat:                 models.Chat{},
		}

		err := repo.Create(questionnaire)
		assert.Nil(t, err)

		questionnaire.Title = "Updated Title"
		err = repo.UpdateWithOwnership(2, questionnaire)
		assert.EqualError(t, err, "user is not the owner of this questionnaire")
	})
	t.Run("TestDeleteWithOwnership", func(t *testing.T) {
		startTime := time.Now().Add(24 * time.Hour)
		endTime := startTime.Add(7 * 24 * time.Hour)
		responseEditDeadline := endTime.Add(24 * time.Hour)

		questionnaire := &models.Questionnaire{
			Title:                "Initial Survey",
			Description:          "Initial Description",
			StartTime:            startTime,
			EndTime:              endTime,
			StepType:             models.StepTypeSequential,
			AllowBacktrack:       true,
			MaxAttempts:          1,
			AnonymityLevel:       models.AnonymityLevelPublic,
			ResponseEditDeadline: responseEditDeadline,
			OwnerID:              1,
			Chat:                 models.Chat{},
		}

		err := repo.Create(questionnaire)
		assert.Nil(t, err)

		err = repo.DeleteWithOwnership(1, questionnaire.ID)
		assert.NoError(t, err)

		result, err := repo.GetByID(questionnaire.ID)
		assert.Nil(t, result)
		assert.EqualError(t, err, "questionnaire not found")
	})
	t.Run("TestDeleteWithOwnership_NotOwner", func(t *testing.T) {
		startTime := time.Now().Add(24 * time.Hour)
		endTime := startTime.Add(7 * 24 * time.Hour)
		responseEditDeadline := endTime.Add(24 * time.Hour)

		questionnaire := &models.Questionnaire{
			Title:                "Initial Survey",
			Description:          "Initial Description",
			StartTime:            startTime,
			EndTime:              endTime,
			StepType:             models.StepTypeSequential,
			AllowBacktrack:       true,
			MaxAttempts:          1,
			AnonymityLevel:       models.AnonymityLevelPublic,
			ResponseEditDeadline: responseEditDeadline,
			OwnerID:              1,
			Chat:                 models.Chat{},
		}

		err := repo.Create(questionnaire)
		assert.Nil(t, err)

		err = repo.DeleteWithOwnership(2, questionnaire.ID)
		assert.EqualError(t, err, "user is not the owner of this questionnaire")
	})

}

func TestValidationQuestionnaire(t *testing.T) {
	db := test_mock.SetupTestDB()
	repo := repositories.NewQuestionnaireRepository(db, false)
	t.Run("TestGetAllActive", func(t *testing.T) {

		activeQuestionnaire := &models.Questionnaire{
			Title:                "Active Survey 2",
			StartTime:            time.Now().Add(-1 * time.Hour),
			EndTime:              time.Now().Add(24 * time.Hour),
			OwnerID:              2,
			StepType:             models.StepTypeRandom,
			AnonymityLevel:       models.AnonymityLevelAnonymous,
			MaxAttempts:          2,
			ResponseEditDeadline: time.Now().Add(24 * time.Hour),
			Chat:                 models.Chat{},
		}
		err := repo.Create(activeQuestionnaire)
		assert.NoError(t, err)

		allActive, err := repo.GetAllActive()
		assert.NoError(t, err)
		assert.Len(t, allActive, 1)
		assert.Equal(t, "Active Survey 2", allActive[0].Title)
	})
	t.Run("TestCreateWithAllFields", func(t *testing.T) {
		startTime := time.Now().Add(24 * time.Hour)
		endTime := startTime.Add(7 * 24 * time.Hour)
		responseEditDeadline := endTime.Add(24 * time.Hour)

		questionnaire := &models.Questionnaire{
			Title:                "Comprehensive Survey",
			Description:          "A detailed survey for testing all fields.",
			StartTime:            startTime,
			EndTime:              endTime,
			StepType:             models.StepTypeSequential,
			AllowBacktrack:       true,
			MaxAttempts:          3,
			AnonymityLevel:       models.AnonymityLevelOwnerOnly,
			ResponseEditDeadline: responseEditDeadline,
			OwnerID:              1,
			Chat:                 models.Chat{},
		}

		err := repo.Create(questionnaire)
		assert.NoError(t, err)
		assert.NotZero(t, questionnaire.ID)

		fetched, err := repo.GetByID(questionnaire.ID)
		assert.NoError(t, err)
		assert.Equal(t, "Comprehensive Survey", fetched.Title)
		assert.Equal(t, "A detailed survey for testing all fields.", fetched.Description)
		assert.Equal(t, startTime.In(time.Local), fetched.StartTime.In(time.Local))
		assert.Equal(t, endTime.In(time.Local), fetched.EndTime.In(time.Local))
		assert.Equal(t, models.StepTypeSequential, fetched.StepType)
		assert.True(t, fetched.AllowBacktrack)
		assert.Equal(t, 3, fetched.MaxAttempts)
		assert.Equal(t, models.AnonymityLevelOwnerOnly, fetched.AnonymityLevel)
		assert.Equal(t, responseEditDeadline.In(time.Local), fetched.ResponseEditDeadline.In(time.Local))
		assert.Equal(t, uint(1), fetched.OwnerID)
	})

	t.Run("TestUpdateWithAllFields", func(t *testing.T) {
		startTime := time.Now().Add(24 * time.Hour)
		endTime := startTime.Add(7 * 24 * time.Hour)
		responseEditDeadline := endTime.Add(24 * time.Hour)

		questionnaire := &models.Questionnaire{
			Title:                "Initial Survey",
			Description:          "Initial Description",
			StartTime:            startTime,
			EndTime:              endTime,
			StepType:             models.StepTypeSequential,
			AllowBacktrack:       true,
			MaxAttempts:          1,
			AnonymityLevel:       models.AnonymityLevelPublic,
			ResponseEditDeadline: responseEditDeadline,
			OwnerID:              1,
			Chat:                 models.Chat{},
		}

		err := repo.Create(questionnaire)
		assert.NoError(t, err)
		assert.NotZero(t, questionnaire.ID)

		newStartTime := startTime.Add(48 * time.Hour)
		newEndTime := newStartTime.Add(14 * 24 * time.Hour)
		newResponseEditDeadline := newEndTime.Add(48 * time.Hour)

		questionnaire.Title = "Updated Survey"
		questionnaire.Description = "Updated Description"
		questionnaire.StartTime = newStartTime
		questionnaire.EndTime = newEndTime
		questionnaire.StepType = models.StepTypeRandom
		questionnaire.AllowBacktrack = false
		questionnaire.MaxAttempts = 5
		questionnaire.AnonymityLevel = models.AnonymityLevelAnonymous
		questionnaire.ResponseEditDeadline = newResponseEditDeadline
		questionnaire.OwnerID = 2

		err = repo.Update(questionnaire)
		assert.NoError(t, err)

		updatedQuestionnaire, err := repo.GetByID(questionnaire.ID)
		assert.NoError(t, err)
		assert.Equal(t, "Updated Survey", updatedQuestionnaire.Title)
		assert.Equal(t, "Updated Description", updatedQuestionnaire.Description)
		assert.Equal(t, newStartTime.In(time.Local), updatedQuestionnaire.StartTime.In(time.Local))
		assert.Equal(t, newEndTime.In(time.Local), updatedQuestionnaire.EndTime.In(time.Local))
		assert.Equal(t, models.StepTypeRandom, updatedQuestionnaire.StepType)
		assert.False(t, updatedQuestionnaire.AllowBacktrack)
		assert.Equal(t, 5, updatedQuestionnaire.MaxAttempts)
		assert.Equal(t, models.AnonymityLevelAnonymous, updatedQuestionnaire.AnonymityLevel)
		assert.Equal(t, newResponseEditDeadline.In(time.Local), updatedQuestionnaire.ResponseEditDeadline.In(time.Local))
		assert.Equal(t, uint(2), updatedQuestionnaire.OwnerID)
	})

	t.Run("TestUpdateWithInvalidData", func(t *testing.T) {
		questionnaire := &models.Questionnaire{
			Title:                "Valid Survey",
			Description:          "Valid Description",
			StartTime:            time.Now().Add(24 * time.Hour),
			EndTime:              time.Now().Add(48 * time.Hour),
			StepType:             models.StepTypeSequential,
			AllowBacktrack:       true,
			MaxAttempts:          2,
			AnonymityLevel:       models.AnonymityLevelPublic,
			ResponseEditDeadline: time.Now().Add(72 * time.Hour),
			OwnerID:              1,
			Chat:                 models.Chat{},
		}

		err := repo.Create(questionnaire)
		assert.NoError(t, err)

		questionnaire.Title = ""
		questionnaire.EndTime = questionnaire.StartTime.Add(-1 * time.Hour)
		questionnaire.StepType = "InvalidType"
		questionnaire.AnonymityLevel = "UnknownLevel"
		questionnaire.MaxAttempts = 0
		questionnaire.OwnerID = 0

		err = repo.Update(questionnaire)
		assert.Error(t, err)
	})

	t.Run("TestCreateWithInvalidData", func(t *testing.T) {

		questionnaire := &models.Questionnaire{
			Description: "Missing title field",
			StartTime:   time.Now(),
			EndTime:     time.Now().Add(24 * time.Hour),
			OwnerID:     1,
		}
		err := repo.Create(questionnaire)
		assert.Error(t, err)

		questionnaire = &models.Questionnaire{
			Title:          "Invalid Time Survey",
			StartTime:      time.Now(),
			EndTime:        time.Now().Add(-24 * time.Hour),
			OwnerID:        1,
			StepType:       models.StepTypeSequential,
			AnonymityLevel: models.AnonymityLevelPublic,
			MaxAttempts:    1,
		}
		err = repo.Create(questionnaire)
		assert.Error(t, err)

		questionnaire = &models.Questionnaire{
			Title:          "Invalid StepType Survey",
			StartTime:      time.Now(),
			EndTime:        time.Now().Add(24 * time.Hour),
			StepType:       "InvalidType",
			AnonymityLevel: models.AnonymityLevelPublic,
			MaxAttempts:    1,
			OwnerID:        1,
		}
		err = repo.Create(questionnaire)
		assert.Error(t, err)

		questionnaire = &models.Questionnaire{
			Title:          "Invalid AnonymityLevel Survey",
			StartTime:      time.Now(),
			EndTime:        time.Now().Add(24 * time.Hour),
			StepType:       models.StepTypeSequential,
			AnonymityLevel: "UnknownLevel",
			MaxAttempts:    1,
			OwnerID:        1,
		}
		err = repo.Create(questionnaire)
		assert.Error(t, err)

		questionnaire = &models.Questionnaire{
			Title:          "Invalid MaxAttempts Survey",
			StartTime:      time.Now(),
			EndTime:        time.Now().Add(24 * time.Hour),
			StepType:       models.StepTypeSequential,
			AnonymityLevel: models.AnonymityLevelPublic,
			MaxAttempts:    0,
		}
		err = repo.Create(questionnaire)
		assert.Error(t, err)

		questionnaire = &models.Questionnaire{
			Title:          "Missing OwnerID Survey",
			StartTime:      time.Now(),
			EndTime:        time.Now().Add(24 * time.Hour),
			StepType:       models.StepTypeSequential,
			AnonymityLevel: models.AnonymityLevelPublic,
			MaxAttempts:    1,
		}
		err = repo.Create(questionnaire)
		assert.Error(t, err)
	})

	t.Run("TestIsOwner", func(t *testing.T) {
		startTime := time.Now().Add(24 * time.Hour)
		endTime := startTime.Add(7 * 24 * time.Hour)
		responseEditDeadline := endTime.Add(24 * time.Hour)

		questionnaire := &models.Questionnaire{
			Title:                "Comprehensive Survey",
			Description:          "A detailed survey for testing all fields.",
			StartTime:            startTime,
			EndTime:              endTime,
			StepType:             models.StepTypeSequential,
			AllowBacktrack:       false,
			MaxAttempts:          3,
			AnonymityLevel:       models.AnonymityLevelOwnerOnly,
			ResponseEditDeadline: responseEditDeadline,
			OwnerID:              1,
			Chat:                 models.Chat{},
		}
		err := repo.Create(questionnaire)
		assert.Nil(t, err)

		isOwner, err := repo.IsOwner(1, questionnaire.ID)
		assert.NoError(t, err)
		assert.True(t, isOwner)

		isOwner, err = repo.IsOwner(2, questionnaire.ID)
		assert.NoError(t, err)
		assert.False(t, isOwner)
	})

	t.Run("TestGetActiveWithOwnership", func(t *testing.T) {
		activeQuestionnaire := &models.Questionnaire{
			Title:                "Active Survey",
			StartTime:            time.Now().Add(-1 * time.Hour),
			EndTime:              time.Now().Add(24 * time.Hour),
			OwnerID:              1,
			StepType:             models.StepTypeSequential,
			AnonymityLevel:       models.AnonymityLevelPublic,
			MaxAttempts:          1,
			ResponseEditDeadline: time.Now().Add(24 * time.Hour),
			Chat:                 models.Chat{},
		}
		err := repo.Create(activeQuestionnaire)
		assert.NoError(t, err)

		inactiveQuestionnaire := &models.Questionnaire{
			Title:                "Inactive Survey",
			StartTime:            time.Now().Add(-48 * time.Hour),
			EndTime:              time.Now().Add(-24 * time.Hour),
			OwnerID:              1,
			StepType:             models.StepTypeSequential,
			AnonymityLevel:       models.AnonymityLevelPublic,
			MaxAttempts:          1,
			ResponseEditDeadline: time.Now().Add(-24 * time.Hour),
			Chat:                 models.Chat{},
		}
		err = repo.Create(inactiveQuestionnaire)
		assert.NoError(t, err)

		activeList, err := repo.GetActiveWithOwnership(1)
		assert.NoError(t, err)
		assert.Len(t, activeList, 1)
		assert.Equal(t, "Active Survey", activeList[0].Title)
	})

}
