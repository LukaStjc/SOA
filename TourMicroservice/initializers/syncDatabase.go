package initializers

import "go-tourm/models"

func SyncDatabase() {
	// DB.AutoMigrate(&models.Tour{})
	DB.AutoMigrate(&models.Tour{}, &models.KeyPoint{}, &models.ShoppingCart{}, &models.OrderItem{}, &models.TourCheckIn{})

}
