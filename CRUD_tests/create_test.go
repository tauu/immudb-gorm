package tests

import (
	"testing"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreate(t *testing.T) {

	// Open connection
	db, err := OpenConnection(t)
	require.NoError(t, err, "There was an error opening connection")

	// Checking if users table exists after creating it
	usersBeforeExists := TableChecker(db, "users")
	assert.Equal(t, false, usersBeforeExists, "Table users should not exists before creating it")

	// Create a users table
	err = db.Debug().AutoMigrate(&User{})
	assert.NoError(t, err, "There was an error creating users table")

	// Checking if users table exists after creating it
	usersAfterExists := TableChecker(db, "users")
	assert.Equal(t, true, usersAfterExists, "Table users should exists after creating it")

	// Define a new record
	uuid, _ := uuid.NewRandom()
	var newUser = User{Name: "Jose", Age: 33, Height: 1.8, CompanyID: uuid}

	// Create a new user record
	res := db.Debug().Create(&newUser)
	assert.NoError(t, res.Error, "An error occurred while creating a new record")

	// Query the data previously inserted to the database
	var user User

	err = db.First(&user, newUser.ID).Error
	if err != nil {
		log.Error().Err(err).Msg("No users found")
	}

	// Test cases
	assert.Equal(t, newUser.Name, user.Name, "The queried name does not match the defined name")
	assert.Equal(t, newUser.Age, user.Age, "The queried age does not match the defined age")
	assert.Equal(t, newUser.Height, user.Height, "The queried height does not match the defined height")
	assert.Equal(t, newUser.CompanyID, user.CompanyID, "The queried companyID does not match the defined companyID")
}
