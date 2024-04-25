// package initializers

// import (
// 	"os"

// 	"gorm.io/driver/postgres"
// 	"gorm.io/gorm"
// )

// var DB *gorm.DB

// func ConnectToDb() {
// 	var err error
// 	dsn := os.Getenv("DB")
// 	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

// 	if err != nil {
// 		panic("Failed to connect to db")
// 	}
// }

package initializers

import (
	"fmt"
	configurations "go-userm/startup"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

var Neo4JDriver neo4j.DriverWithContext

// var Ctx context.Context

func ConnectToDb(config *configurations.Configurations) {
	var err error
	// dsn := os.Getenv("")
	// DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	connectionParams := fmt.Sprintf("user=postgres password=ftn dbname=SOA host=%s port=%s sslmode=disable", config.UserDBHost, config.UserDBPort)
	DB, err = gorm.Open(postgres.Open(connectionParams), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to db")
	}
}

// // ConnectToNeo4j connects to Neo4j using the provided configurations and returns a Neo4j driver.
// func ConnectToNeo4j(config *configurations.Configurations) (neo4j.DriverWithContext, error) {
// 	// Construct the DSN for Neo4j
// 	dsn := fmt.Sprintf("bolt://%s:%s", config.UserGraphDBHost, config.UserGraphDBPort)

// 	// Create a BasicAuth instance using Neo4j username and password from configurations
// 	auth := neo4j.BasicAuth(config.UserGraphDBUsername, config.UserGraphDBPassword, "")

// 	// Create a new Neo4j driver
// 	driver, err := neo4j.NewDriverWithContext(dsn, auth)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to connect to Neo4j: %v", err)
// 	}

// 	return driver, nil
// }
