package pkg

import (
	"QUIC-VPN/internal"
	"log"
	"os"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// EnsureDBFileExists checks if the database file exists and creates it if it doesn't
func EnsureDBFileExists(dbPath string) error {
	dir := filepath.Dir(dbPath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0700)
		if err != nil {
			return err
		}
	}
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		file, err := os.Create(dbPath)
		if err != nil {
			return err
		}
		file.Close()
	}
	return nil
}

// Connect db
func Connect() *gorm.DB {
	dbPath := internal.DbPath

	err := EnsureDBFileExists(dbPath)
	if err != nil {
		log.Fatalln("failed to ensure database file exists: " + err.Error())
	}

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatalln("failed to connect database")
	}
	return db
}

var Db = Connect()

// create table
func CreateTable(table interface{}, db *gorm.DB) {
	err := db.AutoMigrate(table)
	if err != nil {
		log.Fatalln("failed to create table")
	}
}
