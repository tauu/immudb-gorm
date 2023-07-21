package test_migrator

import (
	"net/url"
	"testing"

	"github.com/rs/zerolog/log"
	immudbGorm "github.com/tauu/immudb-gorm"
	_ "github.com/tauu/immusql"
	"gorm.io/gorm"
)

type Employee struct {
	gorm.Model
	Name   string
	Salary int
}

func OpenConnection(t *testing.T) (*gorm.DB, error) {

	// URI to storage location for the database.
	// Example format: immudbe:///folderA/folderB/databaseName
	url := url.URL{
		Scheme: "immudbe",
		Path:   t.TempDir(),
	}

	db, err := gorm.Open(immudbGorm.Open(url.String()), &gorm.Config{})
	if err != nil {
		log.Error().Err(err).Msg("An error occurred while opening connection")
		return nil, err
	}

	return db, nil
}
