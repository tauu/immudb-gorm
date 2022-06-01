package test_associations

import (
	"testing"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestHasMany(t *testing.T) {

	// Open connection
	db, err := OpenConnection()
	if !assert.NoError(t, err, "An error occurred opening connection") {
		t.FailNow()
	}

	// Delete test directory
	defer DeleteTestDir()

	// Check if tables owners and restaurants exists before creating them
	ownersExists := TableChecker("owners", db)
	restaurantsExists := TableChecker("restaurants", db)

	// Test cases
	assert.Equal(t, false, ownersExists, "Table owners already exists before creating it")
	assert.Equal(t, false, restaurantsExists, "Table restaurants already exists before creating it")

	// Create owners table
	err = db.AutoMigrate(&Owner{}, &Restaurant{})
	assert.NoError(t, err, "An error occurred while creating tables")

	// Check if tables owners and restaurants were created
	ownersExists = TableChecker("owners", db)
	restaurantsExists = TableChecker("restaurants", db)

	// Test cases
	assert.Equal(t, true, ownersExists, "There was a problem creating the table owners")
	assert.Equal(t, true, restaurantsExists, "There was a problem creating the table restaurants")

	// Define a new record
	newOwner := Owner{
		Name: "Joel",
		Restaurants: []Restaurant{
			{Name: "Kibuka"},
			{Name: "Carlota Akaneya"},
		},
	}

	// Insert the record previously defined to the data base
	result := db.Create(&newOwner)
	assert.NoError(t, result.Error, "An error occurred inserting a record to the data base")

	var owner Owner
	var restaurants []Restaurant

	ourOwner := Owner{
		Model: gorm.Model{ID: newOwner.ID},
	}

	// Retrieve the restaurants of the owner
	db.Model(&newOwner).Association("Restaurants").Find(&restaurants)
	assert.NoError(t, err, "There was a problem quering the owner")

	// Add another restaurant for the owner in the owner table
	err = db.Model(&ourOwner).Association("Restaurants").Append([]Restaurant{{Name: "Sayonara"}})
	assert.NoError(t, err, "There was a problem adding a restaurant")

	// Retrieve the number of restaurants for the owner
	restaurantsBeforeDel := db.Debug().Model(&newOwner).Association("Restaurants").Count()

	// Remove a restaurant of the owner
	err = db.Model(&newOwner).Association("Restaurants").Delete(&restaurants)
	assert.NoError(t, err, "There was a problem deleting a restaurant")

	// Retrieve the number of restaurants for the owner
	restaurantsAfterDel := db.Debug().Model(&newOwner).Association("Restaurants").Count()

	// Retrieve the data of the owner and it's restaurants
	err = db.Preload("Restaurants").First(&owner, newOwner.ID).Error
	if err != nil {
		log.Error().Err(err).Msg("There was an error quering the first owner")
	}

	// Test cases
	assert.Equal(t, 3, int(restaurantsBeforeDel), "There was an error quering the total amount of restaurants")
	assert.Equal(t, 1, int(restaurantsAfterDel), "There was an error quering the total amount of restaurants")
	assert.Equal(t, 1, int(owner.ID), "There was an error quering owner's ID")
	assert.Equal(t, "Joel", owner.Name, "There was an error quering owner's name")
	assert.Equal(t, "Sayonara", owner.Restaurants[0].Name, "There was an error quering the first restaurant")

}
