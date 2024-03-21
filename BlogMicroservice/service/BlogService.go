package service

import (
	"database-example/model"
	"database-example/repo"
	"fmt"
	"net/http"
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

	url := fmt.Sprintf("http://localhost:3000/is-blocked/%d", blog.UserID)

	// Then make the POST request using the constructed URL
	resp, err1 := http.Post(url, "application/json", nil)
	if err1 != nil {
		return err1
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		return nil
	}

	//user nije blokiran
	err2 := service.BlogRepo.Create(blog)

	if err2 != nil {
		return err2
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
