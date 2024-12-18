package repositories_test

import (
	"testing"

	"github.com/QBG-P2/Voting-System/internal/models"
	"github.com/QBG-P2/Voting-System/internal/repositories"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupVoteTransactionDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	err := db.AutoMigrate(&models.VoteTransaction{}, &models.User{})
	if err != nil {
		return nil
	}
	return db
}

func TestVoteTransactionRepository(t *testing.T) {
	db := setupVoteTransactionDB()
	voteRepo := repositories.NewVoteTransactionRepository(db)

	// Create Users
	fromUser := &models.User{
		NationalID: "1111111111",
		Email:      "from@example.com",
		Password:   "hashedpassword",
	}
	err := db.Create(fromUser).Error
	assert.Nil(t, err)
	assert.NotZero(t, fromUser.ID)

	toUser := &models.User{
		NationalID: "2222222222",
		Email:      "to@example.com",
		Password:   "hashedpassword",
	}
	err = db.Create(toUser).Error
	assert.Nil(t, err)
	assert.NotZero(t, toUser.ID)

	// Create VoteTransaction
	voteTransaction := &models.VoteTransaction{
		BuyerID:  fromUser.ID,
		SellerID: toUser.ID,
		Amount:   50.0,
	}
	err = voteRepo.Create(voteTransaction)
	assert.Nil(t, err)
	assert.NotZero(t, voteTransaction.ID)

	fetchedTransaction, err := voteRepo.GetByID(voteTransaction.ID)
	assert.Nil(t, err)
	assert.NotNil(t, fetchedTransaction)
	assert.Equal(t, float64(50), fetchedTransaction.Amount)

	// Update VoteTransaction
	voteTransaction.Amount = 100.0
	err = voteRepo.Update(voteTransaction)
	assert.Nil(t, err)

	updatedTransaction, err := voteRepo.GetByID(voteTransaction.ID)
	assert.Nil(t, err)
	assert.NotNil(t, updatedTransaction)
	assert.Equal(t, float64(100), updatedTransaction.Amount)

	// Delete VoteTransaction
	err = voteRepo.Delete(voteTransaction.ID)
	assert.Nil(t, err)

	_, err = voteRepo.GetByID(voteTransaction.ID)
	assert.NotNil(t, err)
}
