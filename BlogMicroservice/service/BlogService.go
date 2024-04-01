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

// func (service *BlogService) CreateBlog(blog *model.Blog, authToken string) error {
// 	//DODATI proveru da li je korisnik blokiran
// 	//ako nije, moze da napravi blog
// 	fmt.Printf("Usao u servis")

// 	//url := fmt.Sprintf("http://localhost:3000/is-blocked/%d", blog.UserID)
// 	url := fmt.Sprintf("http://localhost:3000/is-blocked/3")

// 	// Then make the POST request using the constructed URL
// 	resp, err1 := http.Get(url)
// 	fmt.Printf("\nPosle poziva njihovog ms")

// 	if err1 != nil {
// 		fmt.Printf("\nEror jbg")
// 		return err1
// 	}
// 	defer resp.Body.Close()

// 	// Check the response status code
// 	if resp.StatusCode != http.StatusOK {
// 		return nil
// 	}

// 	//user nije blokiran
// 	fmt.Printf("Pre ulaska u repo %s", blog.Title)
// 	err2 := service.BlogRepo.Create(blog)

// 	if err2 != nil {
// 		return err2
// 	}

// 	return nil
// }

func (service *BlogService) CreateBlog(blog *model.Blog, authToken string) error {
	// Construct the URL for checking if the user is blocked
	url := fmt.Sprintf("http://localhost:3000/is-blocked/%d", blog.UserID) // Adjust the URL as needed

	// Create a new HTTP request with the appropriate method, URL, and request body (nil for GET request)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return err
	}

	// Set the Authorization header with the provided auth token
	req.Header.Set("Authorization", "Bearer "+authToken)

	// Send the HTTP request using http.DefaultClient
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("Error sending request: %v\n", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("user is blocked or other error occurred, status code: %d", resp.StatusCode)
	}

	fmt.Printf("Pre ulaska u repo %s", blog.Title)
	err = service.BlogRepo.Create(blog)
	if err != nil {
		fmt.Printf("Error creating blog: %v\n", err)
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
