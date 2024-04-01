package repo

import (
	"database-example/model"

	"gorm.io/gorm"
)

type CommentRepository struct {
	DatabaseConnection *gorm.DB
}

func (repo *CommentRepository) FindById(id string) (model.Comment, error) {
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

func (repo *CommentRepository) FindByBlogId(blogId string) ([]model.Comment, error) {
	var comments []model.Comment
	dbResult := repo.DatabaseConnection.Find(&comments, "blog_id = ?", blogId)

	if dbResult.Error != nil {
		return nil, dbResult.Error
	}

	return comments, nil
}
