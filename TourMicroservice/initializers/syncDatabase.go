package initializers

import "go-tourm/models"

func SyncDatabase() {
	DB.AutoMigrate(&models.Tour{})
}
