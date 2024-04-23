package service

import (
	"database-example/model"
	"database-example/repo"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
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

// func (service *BlogService) CreateBlog(blog *model.Blog) error {
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
// 	err2 := service.BlogRepo.Create(blog)

// 	if err2 != nil {
// 		return err2
// 	}

// 	return nil
// }

func (service *BlogService) CreateComment(comment *model.Comment, authToken string) error {
	fmt.Printf("\nUsao u comment servis")
	// provera da li je korisnik blokiran
	url := fmt.Sprintf("http://localhost:3000/is-blocked/%d", comment.UserID)

	fmt.Printf("User ID received from payload: %d\n", comment.UserID)

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

	fmt.Printf("\nNije blokiran korisnik")

	// provera da li korisnik prati kreatora bloga na kom hoce da ostavi komentar

	//var blog model.Blog
	//blog := find
	//url = fmt.Sprintf("http://localhost:3000/does-follow/%d/%d", comment.UserID, )

	//ako ga prati i nije blokiran, moze da ostavi/napravi komentar

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

	err = service.CommentRepo.Create(comment)

	if err != nil {
		return err
	}

	return nil
}

func (service *BlogService) CreateBlog(blog *model.Blog, authToken string) error {
	// Construct the URL for checking if the user is blocked
	url := fmt.Sprintf("http://user-service:3000/is-blocked/%d", blog.UserID) // Adjust the URL as needed
	fmt.Printf("\nId usera koji pravi blog je %d", blog.UserID)

	fmt.Printf("User ID received from payload: %d\n", blog.UserID)

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

func extractUserIDFromToken(authToken string) (int, error) {
	// Parse the JWT token
	token, err := jwt.Parse(authToken, func(token *jwt.Token) (interface{}, error) {
		// Validate the token signature here if needed
		return []byte("your-secret-key"), nil
	})
	if err != nil {
		return 0, err
	}

	// Extract claims from the token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, fmt.Errorf("invalid token claims")
	}

	// Access user ID from the claims
	userIDFloat, ok := claims["user"].(float64)
	if !ok {
		return 0, fmt.Errorf("user ID not found in token")
	}

	// Convert user ID to int
	userID := int(userIDFloat)

	return userID, nil
}

// In your BlogService
func (service *BlogService) FindAllCommentsByBlogId(blogId uuid.UUID) ([]model.Comment, error) {
	return service.CommentRepo.FindByBlogId(blogId)
}

func (service *BlogService) FindAllBlogs() ([]model.Blog, error) {
	blogs, err := service.BlogRepo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve blogs: %v", err)
	}
	return blogs, nil
}

func (service *BlogService) GetAllComments() ([]model.Comment, error) {
	return service.CommentRepo.GetAllComments()
}
