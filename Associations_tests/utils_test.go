package test_associations

import (
	"net/url"
	"testing"

	"github.com/rs/zerolog/log"
	immudbGorm "github.com/tauu/immudb-gorm"
	_ "github.com/tauu/immusql"
	"gorm.io/gorm"
)

func OpenConnection(t *testing.T) (*gorm.DB, error) {

	// URI to storage location for the database.
	// Example format: immudbe:///folderA/folderB/databaseName
	url := url.URL{
		Scheme: "immudbe",
		Path:   t.TempDir(),
	}

	db, err := gorm.Open(immudbGorm.Open(url.String()), &gorm.Config{
		// Disable foreign keys as immudb is not compatible with them (YET)
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		log.Error().Err(err).Msg("An error occurred while opening connection")
		return nil, err
	}

	return db, nil
}

func OpenImmudbContainer() (*gorm.DB, error) {

	url := url.URL{
		Scheme: "immudb",
		User:   url.UserPassword("immudb", "immudb"),
		Host:   "127.0.0.1:3322",
		Path:   "/defaultdb",
	}
	db, err := gorm.Open(immudbGorm.Open(url.String()), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})

	if err != nil {
		log.Error().Err(err).Msg("An error occurred while opening connection")
		return nil, err
	}

	return db, nil
}

func TableChecker(name string, db *gorm.DB) bool {

	// Pass the name to the function
	result := db.Migrator().HasTable(name)

	return result

}
