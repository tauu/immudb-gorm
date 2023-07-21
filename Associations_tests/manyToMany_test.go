package test_associations

import (
	"fmt"
	"testing"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestManyToMany(t *testing.T) {

	// Alternative variant for testing connection to immudb as a client.
	//db, err := OpenImmudbContainer()

	// Open connection
	db, err := OpenConnection(t)
	require.NoError(t, err, "An error occurred opening connection")

	// Check if tables students and languages exists before creating them
	studentsExists := TableChecker("students", db)
	languagesExists := TableChecker("languages", db)

	// Test cases
	assert.Equal(t, false, studentsExists, "Table students already exists before creating it")
	assert.Equal(t, false, languagesExists, "Table languages already exists before creating it")

	// Create students table
	err = db.Debug().AutoMigrate(&Student{})
	assert.NoError(t, err, "An error occurred while creating tables")

	// Check if tables students and languages were created
	studentsExists = TableChecker("students", db)
	languagesExists = TableChecker("languages", db)

	// Test cases
	assert.Equal(t, true, studentsExists, "There was a problem creating the table students")
	assert.Equal(t, true, languagesExists, "There was a problem creating the table languages")

	// Define a new record
	newStudent := Student{
		Name: "Joel",
		Age:  32,
		Languages: []Language{
			{Name: "Spanish"},
			{Name: "Catalan"},
			{Name: "English"},
			{Name: "German"},
		},
	}

	//creates a row in the student table and in the languages table.
	result := db.Debug().Create(&newStudent)
	assert.NoError(t, result.Error, "There was an error inserting a new record")

	var student Student
	var languages []Language

	ourStudent := Student{
		Model: gorm.Model{ID: newStudent.ID},
	}

	// Retrieve the languages of the student
	err = db.Debug().Model(&newStudent).Association("Languages").Find(&languages)
	assert.NoError(t, err, "There was an error retrieving the languages of the student")

	// Add another language for the student in the student table
	err = db.Debug().Model(&ourStudent).Association("Languages").Append([]Language{{Name: "Russian"}})
	assert.NoError(t, err, "There was an error adding a language")

	// Retrieve the number of languages for the student
	languagesBeforeDel := db.Debug().Model(&newStudent).Association("Languages").Count()
	// Returns 0. Not working
	fmt.Println("before delete", languagesBeforeDel)

	// Removes the first 4 languages of the student
	err = db.Model(&newStudent).Association("Languages").Delete(&languages)
	assert.NoError(t, err, "There was an error deleting a language")

	// Retrieve the number of languages for the student
	languagesAfterDel := db.Debug().Model(&newStudent).Association("Languages").Count()
	// Returns 0. Not working
	fmt.Println("after delete", languagesAfterDel)

	// Retrieve the data of the student and it's languages
	err = db.Preload("Languages").First(&student, &newStudent.ID).Error
	if err != nil {
		log.Error().Err(err).Msg("There was an error querying the first student")
	}

	// Test cases
	assert.Equal(t, 5, languagesBeforeDel, "An error occurred parsing languages")
	assert.Equal(t, 1, languagesAfterDel, "An error occurred parsing languages")
	assert.Equal(t, "Joel", student.Name, "An error occurred parsing the languages")
	assert.Equal(t, 32, student.Age, "An error occurred parsing the languages")
	assert.Equal(t, "Russian", student.Languages[0].Name, "An error occurred quering the first langauge")

}
