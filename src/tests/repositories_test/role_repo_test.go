package repositories_test

import (
	"testing"

	"github.com/QBG-P2/Voting-System/internal/models"
	"github.com/QBG-P2/Voting-System/internal/repositories"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupRoleDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	err := db.AutoMigrate(&models.Role{}, &models.Permission{}, &models.User{})
	if err != nil {
		return nil
	}
	return db
}

func TestRoleRepository(t *testing.T) {
	db := setupRoleDB()
	roleRepo := repositories.NewRoleRepository(db)

	role := &models.Role{
		Name: "Admin",
	}
	err := roleRepo.Create(role)
	assert.Nil(t, err)
	assert.NotZero(t, role.ID)

	duplicateRole := &models.Role{
		Name: "Admin",
	}
	err = roleRepo.Create(duplicateRole)
	assert.NotNil(t, err)

	fetchedRole, err := roleRepo.GetByID(role.ID)
	assert.Nil(t, err)
	assert.NotNil(t, fetchedRole)
	assert.Equal(t, "Admin", fetchedRole.Name)

	role.Name = "SuperAdmin"
	err = roleRepo.Update(role)
	assert.Nil(t, err)

	updatedRole, err := roleRepo.GetByID(role.ID)
	assert.Nil(t, err)
	assert.NotNil(t, updatedRole)
	assert.Equal(t, "SuperAdmin", updatedRole.Name)

	err = roleRepo.Delete(role.ID)
	assert.Nil(t, err)

	_, err = roleRepo.GetByID(role.ID)
	assert.NotNil(t, err)
}
