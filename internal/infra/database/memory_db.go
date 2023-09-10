package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
)

func prepareDB(t *testing.T) (*gorm.DB, func()) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}

	cleanup := createCleanupFunction(t, db)

	return db, cleanup
}

func createCleanupFunction(t *testing.T, db *gorm.DB) func() {
	return func() {
		if sqlDB, err := db.DB(); err == nil {
			if err := sqlDB.Close(); err != nil {
				t.Fatal(err)
			}
		}
	}
}

func migrateTable(db *gorm.DB, t *testing.T, model interface{}) {
	err := db.AutoMigrate(model)
	if err != nil {
		t.Fatal(err)
	}
}

func initDBTest(t *testing.T, model interface{}) (*gorm.DB, func()) {
	db, cleanup := prepareDB(t)
	migrateTable(db, t, model)
	return db, cleanup
}
