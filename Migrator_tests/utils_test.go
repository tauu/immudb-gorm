package test_migrator

import (
	"net/url"
	"os"

	"github.com/rs/zerolog/log"
	immudbGorm "github.com/tauu/immudb-gorm"
	"gorm.io/gorm"
)

type Employee struct {
	gorm.Model
	Name   string
	Salary int
}

func OpenConnection() (*gorm.DB, error) {

	// Creates a database object
	// http://foo/asdfadf/test
	// "file:///test.html"
	path, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	// URI to storage location for the database.
	// Example format: immudbe:///folderA/folderB/databaseName
	url := url.URL{
		Scheme: "immudbe",
		Path:   path + "/test/testdb",
	}

	db, err := gorm.Open(immudbGorm.Open(url.String()), &gorm.Config{})
	if err != nil {
		log.Error().Err(err).Msg("An error occurred while opening connection")
		return nil, err
	}

	return db, nil
}

func DeleteTestDir() {

	// Deletes the test directory
	removeDirError := os.RemoveAll("test")
	if removeDirError != nil {
		log.Error().Err(removeDirError).Msg("An error occurred while deleting test directory")
	}
}
