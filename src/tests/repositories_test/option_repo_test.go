package repositories_test

import (
	"testing"

	"github.com/QBG-P2/Voting-System/internal/models"
	"github.com/QBG-P2/Voting-System/internal/repositories"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupOptionDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	err := db.AutoMigrate(&models.Option{}, &models.Question{})
	if err != nil {
		return nil
	}
	return db
}

func TestOptionRepository(t *testing.T) {
	db := setupOptionDB()
	optionRepo := repositories.NewOptionRepository(db)

	question := &models.Question{
		QuestionnaireID: 1,
	}
	err := db.Create(question).Error
	assert.Nil(t, err)
	assert.NotZero(t, question.ID)

	option := &models.Option{
		Text:       "Option 1",
		QuestionID: question.ID,
	}
	err = optionRepo.Create(option)
	assert.Nil(t, err)
	assert.NotZero(t, option.ID)

	fetchedOption, err := optionRepo.GetByID(option.ID)
	assert.Nil(t, err)
	assert.NotNil(t, fetchedOption)
	assert.Equal(t, "Option 1", fetchedOption.Text)

	var options []models.Option
	err = db.Model(&models.Option{}).Where("question_id = ?", question.ID).Find(&options).Error
	assert.Nil(t, err)
	assert.Len(t, options, 1)

	var total int64
	err = db.Model(&models.Option{}).Where("question_id = ?", question.ID).Count(&total).Error
	assert.Nil(t, err)
	assert.Equal(t, int64(1), total)

	option.Text = "Updated Option"
	err = optionRepo.Update(option)
	assert.Nil(t, err)

	updatedOption, err := optionRepo.GetByID(option.ID)
	assert.Nil(t, err)
	assert.NotNil(t, updatedOption)
	assert.Equal(t, "Updated Option", updatedOption.Text)

	err = optionRepo.Delete(option.ID)
	assert.Nil(t, err)

	_, err = optionRepo.GetByID(option.ID)
	assert.NotNil(t, err)
}
