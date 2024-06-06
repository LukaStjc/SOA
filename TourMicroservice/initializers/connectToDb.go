package initializers

import (
	"fmt"
	configurations "go-tourm/startup"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDb(config *configurations.Configurations) {
	var err error
	// dsn := os.Getenv("")
	// DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	connectionParams := fmt.Sprintf("user=postgres password=ftn dbname=SOA host=%s port=%s sslmode=disable", config.TourDBHost, config.TourDBPort)

	DB, err = gorm.Open(postgres.Open(connectionParams), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to db")
	}

	DB = AddQueryHook(DB)
}

// AddQueryHook adds the query hook to the provided *gorm.DB instance
func AddQueryHook(db *gorm.DB) *gorm.DB {
	// Create a new hook
	hook := &QueryHook{}

	// Add the hook to the *gorm.DB instance
	db.Callback().Query().Before("gorm:query").Register("QueryHook", hook.BeforeQuery)
	db.Callback().Query().After("gorm:query").Register("QueryHook", hook.AfterQuery)

	return db
}

// QueryHook is a custom query hook for GORM
type QueryHook struct{}

// BeforeQuery is called before executing each query
func (q *QueryHook) BeforeQuery(db *gorm.DB) {
	// Add code to be executed before each query
	// For example, you can start a span here
}

// AfterQuery is called after executing each query
func (q *QueryHook) AfterQuery(db *gorm.DB) {
	// Add code to be executed after each query
	// For example, you can end the span here
}
