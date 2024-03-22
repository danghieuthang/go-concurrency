package concurrency

import (
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"testing"
)

// Db is a global variable that holds the database connection.
var Db *gorm.DB

// TestEntity is a struct that represents a table in the database.
type TestEntity struct {
	ID      uint
	Name    string
	Version Version
}

// init is a function that is automatically executed when the package is imported.
// It opens a database connection and migrates the TestEntity table.
func init() {
	var err error
	if Db, err = OpenTestConnection(); err != nil {
		log.Printf("Error opening test connection: %v", err)
		os.Exit(1)
	}

	sqlDb, err := Db.DB()
	if err == nil {
		err = sqlDb.Ping()
	}
	if err != nil {
		log.Printf("Error pinging test connection: %v", err)
	}

	// Drop all tables if they exist
	if err := Db.Migrator().DropTable(&TestEntity{}); err != nil {
		log.Printf("Error dropping tables: %v", err)
	}

	// Migrate the tables
	Db.AutoMigrate(&TestEntity{})
}

// OpenTestConnection is a function that opens a database connection.
func OpenTestConnection() (db *gorm.DB, err error) {
	dns := "host=localhost user=postgres password=1 dbname=concurrency_test port=5432"
	db, err = gorm.Open(postgres.Open(dns), &gorm.Config{})
	return db.Debug(), err
}

// TestAutoSetIfEmpty is a test function that checks if the Version field is automatically set when a new TestEntity is created.
func TestAutoSetIfEmpty(t *testing.T) {
	e := TestEntity{
		ID:   1,
		Name: "insert new",
	}
	err := Db.Create(&e).Error
	assert.NoError(t, err)
	assert.True(t, e.Version.Valid)
	assert.NotEmpty(t, e.Version.String)
}

// TestConcurrency is a test function that checks if the Version field is updated when a TestEntity is updated,
// and if an update fails when an outdated Version is used.
func TestConcurrency(t *testing.T) {
	e := TestEntity{
		ID:   3,
		Name: "created",
	}
	err := Db.Create(&e).Error
	assert.NoError(t, err)

	var ec TestEntity
	err = Db.First(&ec, "id", 3).Error
	assert.NoError(t, err)
	assert.Equal(t, e.ID, ec.ID)

	tx := Db.Model(&e).Update("name", "first name")
	assert.NoError(t, tx.Error)
	assert.Equal(t, int64(1), tx.RowsAffected)
	assert.Equal(t, e.Name, "first name")

	assert.True(t, e.Version.Valid)
	assert.NotEmpty(t, e.Version.String)
	assert.NotEqual(t, e.Version.String, ec.Version.String)

	affected := Db.Model(&ec).Update("name", "second time").RowsAffected
	assert.Equal(t, int64(0), affected)
}
