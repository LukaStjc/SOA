package repo

import (
	"database-example/model"
	"fmt"

	"gorm.io/gorm"
)

type BlogRepository struct {
	DatabaseConnection *gorm.DB
}

func (repo *BlogRepository) FindById(id string) (model.Blog, error) {
	blog := model.Blog{}
	dbResult := repo.DatabaseConnection.First(&blog, "id = ?", id)

	if dbResult != nil {
		return blog, dbResult.Error
	}

	return blog, nil
}

func (repo *BlogRepository) Create(blog *model.Blog) error {

	fmt.Printf("pocetak kreiranja %s", blog.Title)

	dbResult := repo.DatabaseConnection.Create(blog)

	fmt.Println("Kreiran")

	if dbResult.Error != nil {
		return dbResult.Error
	}

	println("Rows affected: ", dbResult.RowsAffected)
	return nil
}

func (repo *BlogRepository) FindAll() ([]model.Blog, error) {
	var blogs []model.Blog
	dbResult := repo.DatabaseConnection.Find(&blogs)

	if dbResult != nil {
		return blogs, dbResult.Error
	}

	return blogs, nil
}
