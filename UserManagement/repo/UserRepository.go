package repo

import (
	"database-example/model"

	"gorm.io/gorm"
)

type UserRepository struct {
	DatabaseConnection *gorm.DB
}

// func (repo *UserRepository) FindById(id string) (model.Student, error) {
// 	student := model.Student{}
// 	dbResult := repo.DatabaseConnection.First(&student, "id = ?", id)
// 	if dbResult != nil {
// 		return student, dbResult.Error
// 	}
// 	return student, nil
// }

func (repo *UserRepository) CreateUser(user *model.User) error {
	dbResult := repo.DatabaseConnection.Create(user)
	if dbResult.Error != nil {
		return dbResult.Error
	}
	println("Rows affected: ", dbResult.RowsAffected)
	return nil
}
