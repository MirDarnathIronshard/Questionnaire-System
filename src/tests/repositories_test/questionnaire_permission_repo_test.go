package repositories_test

import (
	"testing"

	"github.com/QBG-P2/Voting-System/internal/models"
	"github.com/QBG-P2/Voting-System/internal/repositories"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupQuestionnairePermissionDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	err := db.AutoMigrate(&models.QuestionnairePermission{}, &models.QuestionnaireRole{})
	if err != nil {
		return nil
	}
	return db
}

func TestQuestionnairePermissionRepository(t *testing.T) {
	db := setupQuestionnairePermissionDB()
	permissionRepo := repositories.NewQuestionnairePermissionRepository(db)

	permission := &models.QuestionnairePermission{
		Name: "ViewResults",
	}
	err := permissionRepo.Create(permission)
	assert.Nil(t, err)
	assert.NotZero(t, permission.ID)

	duplicatePermission := &models.QuestionnairePermission{
		Name: "ViewResults",
	}
	err = permissionRepo.Create(duplicatePermission)
	assert.NotNil(t, err)

	fetchedPermission, err := permissionRepo.GetByID(permission.ID)
	assert.Nil(t, err)
	assert.NotNil(t, fetchedPermission)
	assert.Equal(t, "ViewResults", fetchedPermission.Name)

	permission.Name = "EditResults"
	err = permissionRepo.Update(permission)
	assert.Nil(t, err)

	updatedPermission, err := permissionRepo.GetByID(permission.ID)
	assert.Nil(t, err)
	assert.NotNil(t, updatedPermission)
	assert.Equal(t, "EditResults", updatedPermission.Name)

	err = permissionRepo.Delete(permission.ID)
	assert.Nil(t, err)

	_, err = permissionRepo.GetByID(permission.ID)
	assert.NotNil(t, err)
}
