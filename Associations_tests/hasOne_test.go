package test_associations

import (
	"testing"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHasOne(t *testing.T) {

	// Open connection
	db, err := OpenConnection(t)
	require.NoError(t, err, "An error occurred opening connection")

	// Check if tables users and credit_cards exists before creating them
	usersExists := TableChecker("users", db)
	creditCardsExists := TableChecker("credit_cards", db)

	// Test cases
	assert.Equal(t, false, usersExists, "Table users already exists before creating it")
	assert.Equal(t, false, creditCardsExists, "Table credit_cards already exists before creating it")

	// Create users table
	err = db.AutoMigrate(&User{}, &CreditCard{})
	assert.NoError(t, err, "An error occurred while creating tables")

	// Check if tables users and credit_cards were created
	usersExists = TableChecker("users", db)
	creditCardsExists = TableChecker("credit_cards", db)

	// Test cases
	assert.Equal(t, true, usersExists, "There was a problem creating the table users")
	assert.Equal(t, true, creditCardsExists, "There was a problem creating the table credit_cards")

	// Define a new record
	newUser := User{
		Name: "Joel",
		CreditCard: CreditCard{
			Number: "111122223333",
		},
	}

	// Insert a new record to the data base
	result := db.Create(&newUser)
	assert.NoError(t, result.Error, "There was an error inserting a new record")

	var user User
	var creditCard CreditCard

	// ourUser := Owner{
	// 	Model: gorm.Model{ID: newUser.ID},
	// }

	// Retrieve the credit cards of the user
	err = db.Model(&newUser).Association("CreditCard").Find(&creditCard)
	assert.NoError(t, err, "There was an error retrieving the credit cards of the user")

	// The following test is deabled until the next version of immudb come out
	// Then .Replace() method will work without any issues.

	// // Replace the credit card of the user
	// err = db.Debug().Model(&newUser).Association("CreditCard").Replace(&CreditCard{Number: "0000000"})
	// assert.NoError(t, err, "There was an error replacing the credit card of the user")

	// HERE will come a test case to check that the credit card number was updated correctly

	// Retrieve the number of credit cards for the user
	creditCardsBeforeDel := db.Model(&newUser).Association("CreditCard").Count()

	// Remove a credit card of the user
	err = db.Model(&newUser).Association("CreditCard").Delete(&creditCard)
	assert.NoError(t, err, "There was an error deleting the credit card of the user")

	// Retrieve the number of credit cards for the user
	creditCardsAfterDel := db.Model(&newUser).Association("CreditCard").Count()

	// Retrieve the data of the user and it's credit card
	err = db.Preload("CreditCard").First(&user, &newUser.ID).Error
	if err != nil {
		log.Error().Err(err).Msg("There was an error quering the first user")
	}

	// Test cases
	assert.Equal(t, 1, int(creditCardsBeforeDel), "There was an error quering the total amount of credit cards")
	assert.Equal(t, 0, int(creditCardsAfterDel), "There was an error quering the total amount of credit cards")
	assert.Equal(t, 1, int(user.ID), "There was an error quering user's ID")
	assert.Equal(t, "Joel", user.Name, "There was an error quering user's name")
	assert.Equal(t, "", user.CreditCard.Number, "There was an error quering the number of the credit card")

}
