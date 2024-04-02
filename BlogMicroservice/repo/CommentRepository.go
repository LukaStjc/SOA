package repo

import (
	"database-example/model"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CommentRepository struct {
	DatabaseConnection *gorm.DB
}

func (repo *CommentRepository) FindById(id string) (model.Comment, error) {

	if repo.DatabaseConnection == nil {
		fmt.Println("DatabaseConnection is nil")
	} else {
		fmt.Println("DatabaseConnection is initialized")
	}

	comment := model.Comment{}
	dbResult := repo.DatabaseConnection.First(&comment, "id = ?", id)

	if dbResult != nil {
		return comment, dbResult.Error
	}

	return comment, nil
}

func (repo *CommentRepository) Create(comment *model.Comment) error {
	dbResult := repo.DatabaseConnection.Create(comment)

	if dbResult.Error != nil {
		return dbResult.Error
	}

	println("Rows affected: ", dbResult.RowsAffected)
	return nil
}

func (repo *CommentRepository) FindByBlogId(blogId uuid.UUID) ([]model.Comment, error) {
	var comments []model.Comment
	// Use GORM's Where method with the correct type for blogId
	err := repo.DatabaseConnection.Where("blog_id = ?", blogId).Find(&comments).Error
	if err != nil {
		return nil, err
	}
	return comments, nil
}

func (repo *CommentRepository) GetAllComments() ([]model.Comment, error) {
	var comments []model.Comment
	result := repo.DatabaseConnection.Find(&comments)
	return comments, result.Error
}
