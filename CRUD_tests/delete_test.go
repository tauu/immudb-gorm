package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeleteSoft(t *testing.T) {

	// Open connection
	db, err := OpenConnection()
	if !assert.NoError(t, err, "There was an error openning connection") {
		t.FailNow()
	}

	// Delete the test directory
	defer DeleteTestDir()

	// Checking if users table exists after creating it
	usersBeforeExists := TableChecker(db, "users")
	assert.Equal(t, false, usersBeforeExists, "Table users should not exists before creating it")

	// Create a users table
	err = db.AutoMigrate(&User{})
	assert.NoError(t, err, "There was an error creating users table")

	// Checking if users table exists after creating it
	usersAfterExists := TableChecker(db, "users")
	assert.Equal(t, true, usersAfterExists, "Table users should exists after creating it")

	// Define a new record
	var newUser = User{Name: "Jose", Age: 33}

	// Create a new user record
	res := db.Create(&newUser)
	assert.NoError(t, res.Error, "An error occurred while creating a new record")

	// Delete a specific record softly, matching the parameters introduced
	result := db.Where("id = ?", newUser.ID).Delete(&User{})
	assert.NoError(t, result.Error, "An error occurred while deleting")

	// Check that the user was deleted
	var user User
	err = db.First(&user, newUser.ID).Error
	assert.Equal(t, err.Error(), "record not found", "There was an error Deleting the row from the database")

}
