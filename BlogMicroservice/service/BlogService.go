package service

import (
	"database-example/model"
	"database-example/repo"
	"fmt"
)

type BlogService struct {
	BlogRepo    *repo.BlogRepository
	CommentRepo *repo.CommentRepository
}

func (service *BlogService) FindBlogById(id string) (*model.Blog, error) {
	blog, err := service.BlogRepo.FindById(id)

	if err != nil {
		return nil, fmt.Errorf(fmt.Sprintf("menu item with id %s not found", id))
	}

	return &blog, nil
}

func (service *BlogService) FindCommentById(id string) (*model.Comment, error) {
	comment, err := service.CommentRepo.FindById(id)

	if err != nil {
		return nil, fmt.Errorf(fmt.Sprintf("menu item with id %s not found", id))
	}

	return &comment, nil
}

func (service *BlogService) CreateBlog(blog *model.Blog) error {
	//DODATI proveru da li je korisnik blokiran
	//ako nije, moze da napravi blog

	err := service.BlogRepo.Create(blog)

	if err != nil {
		return err
	}

	return nil
}

func (service *BlogService) CreateComment(comment *model.Comment) error {
	//DODATI proveru da li je korisnik blokiran
	//DODATI proveru da li korisnik prati kreatora bloga na kom hoce da ostavi komentar
	//ako ga prati i nije blokiran, moze da ostavi/napravi komentar

	err := service.CommentRepo.Create(comment)

	if err != nil {
		return err
	}

	return nil
}
