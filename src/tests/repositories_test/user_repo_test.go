package repositories_test

import (
	"testing"

	"github.com/QBG-P2/Voting-System/internal/models"
	"github.com/QBG-P2/Voting-System/internal/repositories"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupUserDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	err := db.AutoMigrate(&models.User{}, &models.Role{})
	if err != nil {
		return nil
	}
	return db
}

func TestUserRepository(t *testing.T) {
	db := setupUserDB()
	userRepo := repositories.NewUserRepository(db)

	// Create User
	user := &models.User{
		NationalID: "1234567890",
		Email:      "user@example.com",
		Password:   "hashedpassword",
	}
	err := userRepo.CreateUser(user)
	assert.Nil(t, err)
	assert.NotZero(t, user.ID)

	// Attempt to create duplicate User with same NationalID
	duplicateUser := &models.User{
		NationalID: "1234567890",
		Email:      "duplicate@example.com",
		Password:   "hashedpassword",
	}
	err = userRepo.CreateUser(duplicateUser)
	assert.NotNil(t, err)

	// Get User by ID
	fetchedUser, err := userRepo.GetByID(user.ID)
	assert.Nil(t, err)
	assert.NotNil(t, fetchedUser)
	assert.Equal(t, "user@example.com", fetchedUser.Email)

	// Update User
	user.Email = "updated@example.com"
	err = userRepo.Update(user)
	assert.Nil(t, err)

	updatedUser, err := userRepo.GetByID(user.ID)
	assert.Nil(t, err)
	assert.NotNil(t, updatedUser)
	assert.Equal(t, "updated@example.com", updatedUser.Email)

	// Delete User
	err = userRepo.Delete(user.ID)
	assert.Nil(t, err)

	_, err = userRepo.GetByID(user.ID)
	assert.NotNil(t, err)
}
