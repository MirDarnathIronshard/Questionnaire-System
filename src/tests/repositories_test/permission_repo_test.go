package repositories_test

import (
	"testing"

	"github.com/QBG-P2/Voting-System/internal/models"
	"github.com/QBG-P2/Voting-System/internal/repositories"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupPermissionDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	err := db.AutoMigrate(&models.Permission{}, &models.Role{})
	if err != nil {
		return nil
	}
	return db
}

func TestPermissionRepository(t *testing.T) {
	db := setupPermissionDB()
	permissionRepo := repositories.NewPermissionRepository(db)

	permission := &models.Permission{
		Name:   "ReadPermission",
		Path:   "/api/resource",
		Method: "GET",
	}
	err := permissionRepo.Create(permission)
	assert.Nil(t, err)
	assert.NotZero(t, permission.ID)

	duplicatePermission := &models.Permission{
		Name:   "ReadPermission",
		Path:   "/api/duplicate",
		Method: "GET",
	}
	err = permissionRepo.Create(duplicatePermission)
	assert.NotNil(t, err)

	fetchedPermission, err := permissionRepo.GetByID(permission.ID)
	assert.Nil(t, err)
	assert.NotNil(t, fetchedPermission)
	assert.Equal(t, "ReadPermission", fetchedPermission.Name)

	permission.Path = "/api/updated"
	err = permissionRepo.Update(permission)
	assert.Nil(t, err)

	updatedPermission, err := permissionRepo.GetByID(permission.ID)
	assert.Nil(t, err)
	assert.NotNil(t, updatedPermission)
	assert.Equal(t, "/api/updated", updatedPermission.Path)

	err = permissionRepo.Delete(permission.ID)
	assert.Nil(t, err)

	_, err = permissionRepo.GetByID(permission.ID)
	assert.NotNil(t, err)
}
