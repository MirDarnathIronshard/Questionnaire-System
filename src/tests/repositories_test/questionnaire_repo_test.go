package repositories_test

import (
	"testing"
	"time"

	"github.com/QBG-P2/Voting-System/internal/models"
	"github.com/QBG-P2/Voting-System/internal/repositories"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupQuestionnaireDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	err := db.AutoMigrate(
		&models.Questionnaire{},
		&models.User{},
		&models.Question{},
		&models.Chat{},
		&models.QuestionnaireRole{},
		&models.UserQuestionnaireAccess{},
	)
	if err != nil {
		return nil
	}
	return db
}

func TestQuestionnaireRepository(t *testing.T) {
	db := setupQuestionnaireDB()
	questionnaireRepo := repositories.NewQuestionnaireRepository(db, false)

	// Create User as Owner
	user := &models.User{
		NationalID: "1234567890",
		Email:      "owner@example.com",
		Password:   "hashedpassword",
	}
	err := db.Create(user).Error
	assert.Nil(t, err)
	assert.NotZero(t, user.ID)

	// Create Questionnaire
	startTime := time.Now()
	endTime := startTime.Add(24 * time.Hour)
	questionnaire := &models.Questionnaire{
		Title:          "Sample Questionnaire",
		OwnerID:        user.ID,
		StartTime:      startTime,
		EndTime:        endTime,
		AnonymityLevel: "Public",
		StepType:       "Sequential",
	}
	err = questionnaire.Validate()
	assert.Nil(t, err)

	err = questionnaireRepo.Create(questionnaire)
	assert.Nil(t, err)
	assert.NotZero(t, questionnaire.ID)

	// Add User Access to Questionnaire
	access := &models.UserQuestionnaireAccess{
		UserID:          user.ID,
		QuestionnaireID: questionnaire.ID,
	}
	err = db.Create(access).Error
	assert.Nil(t, err)
	assert.NotZero(t, access.ID)

	// Get Questionnaire by ID
	fetchedQuestionnaire, err := questionnaireRepo.GetByID(questionnaire.ID)
	assert.Nil(t, err)
	assert.NotNil(t, fetchedQuestionnaire)
	assert.Equal(t, "Sample Questionnaire", fetchedQuestionnaire.Title)

	// Update Questionnaire
	questionnaire.Title = "Updated Questionnaire"
	err = questionnaire.Validate()
	assert.Nil(t, err)

	err = questionnaireRepo.Update(questionnaire)
	assert.Nil(t, err)

	updatedQuestionnaire, err := questionnaireRepo.GetByID(questionnaire.ID)
	assert.Nil(t, err)
	assert.NotNil(t, updatedQuestionnaire)
	assert.Equal(t, "Updated Questionnaire", updatedQuestionnaire.Title)

	// Delete Questionnaire
	err = questionnaireRepo.Delete(questionnaire.ID)
	assert.Nil(t, err)

	_, err = questionnaireRepo.GetByID(questionnaire.ID)
	assert.NotNil(t, err)
}
