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
	// fmt.Printf("Usao u servis")

	// url := fmt.Sprintf("http://localhost:3000/is-blocked/%d", blog.UserID)

	// // Then make the POST request using the constructed URL
	// resp, err1 := http.Get(url)
	// fmt.Printf("\nPosle poziva njihovog ms")

	// if err1 != nil {
	// 	fmt.Printf("\nEror neki")
	// 	return err1
	// }
	// defer resp.Body.Close()

	// // Check the response status code
	// if resp.StatusCode != http.StatusOK {
	// 	return nil
	// }

	// //user nije blokiran
	// fmt.Printf("Pre ulaska u repo %s", blog.Title)
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

	// url := fmt.Sprintf("http://localhost:3000/is-blocked/%d", comment.UserID)

	// // Then make the GET request using the constructed URL
	// resp, err1 := http.Get(url)

	// if err1 != nil {
	// 	return err1
	// }
	// defer resp.Body.Close()

	// // Check the response status code
	// if resp.StatusCode != http.StatusOK {
	// 	return nil
	// }

	// //user nije blokiran

	// blog, _ := service.BlogRepo.FindById(comment.BlogID.String())

	// url := fmt.Sprintf("http://localhost:3000/does-follow/%d/%d", comment.UserID, blog.UserID)

	// // Then make the GET request using the constructed URL
	// resp1, err3 := http.Get(url)

	// if err3 != nil {
	// 	return err3
	// }
	// defer resp1.Body.Close()

	// // Check the response status code
	// if resp1.StatusCode != http.StatusOK {
	// 	return nil
	// }

	// comment creator follows blog creator
	err := service.CommentRepo.Create(comment)

	if err != nil {
		return err
	}

	return nil
}
