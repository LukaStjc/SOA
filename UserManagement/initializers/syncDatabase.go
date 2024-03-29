package initializers

import "go-userm/models"

func SyncDatabase() {
	DB.AutoMigrate(&models.User{})
}
