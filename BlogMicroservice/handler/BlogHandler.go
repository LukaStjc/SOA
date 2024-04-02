package handler

import (
	"database-example/model"
	"database-example/service"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type BlogHandler struct {
	BlogService *service.BlogService
}

func (handler *BlogHandler) GetBlog(writer http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	log.Printf("Blog sa id-em %s", id)
	blog, err := handler.BlogService.FindBlogById(id)
	writer.Header().Set("Content-Type", "application/json")

	if err != nil {
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(blog)
}

func (handler *BlogHandler) GetComment(writer http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	log.Printf("Komentar sa id-em %s", id)
	comment, err := handler.BlogService.FindCommentById(id)
	writer.Header().Set("Content-Type", "application/json")

	if err != nil {
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(comment)
}

func (handler *BlogHandler) CreateBlog(writer http.ResponseWriter, req *http.Request) {
	var blog model.Blog
	err := json.NewDecoder(req.Body).Decode(&blog)

	if err != nil {
		println("Error while parsing json")
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	// TOKEN
	authHeader := req.Header.Get("Authorization")
	if authHeader == "" {
		println("Missing Authorization header")
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}
	authToken := authHeader[len("Bearer "):]
	fmt.Println("Auth Token:", authToken)

	err = handler.BlogService.CreateBlog(&blog, authToken)

	if err != nil {
		println("Error while creating a new blog")
		writer.WriteHeader(http.StatusExpectationFailed)
		return
	}

	writer.WriteHeader(http.StatusCreated)
	writer.Header().Set("Content-Type", "application/json")
}

func (handler *BlogHandler) CreateComment(writer http.ResponseWriter, req *http.Request) {
	var comment model.Comment
	err := json.NewDecoder(req.Body).Decode(&comment)

	if err != nil {
		println("Error while parsing json")
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	err = handler.BlogService.CreateComment(&comment)

	if err != nil {
		println("Error while creating a new comment")
		writer.WriteHeader(http.StatusExpectationFailed)
		return
	}

	writer.WriteHeader(http.StatusCreated)
	writer.Header().Set("Content-Type", "application/json")
}

func (handler *BlogHandler) GetAllCommentsByBlogId(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	blogIdStr := vars["blogId"] // Ensure you're using "github.com/gorilla/mux" and have a route variable named "blogId"
	blogId, err := uuid.Parse(blogIdStr)
	if err != nil {
		http.Error(w, "Invalid Blog ID format", http.StatusBadRequest)
		return
	}

	comments, err := handler.BlogService.FindAllCommentsByBlogId(blogId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comments)
}

func (handler *BlogHandler) FindAllBlogs(writer http.ResponseWriter, req *http.Request) {
	blogs, err := handler.BlogService.FindAllBlogs()
	writer.Header().Set("Content-Type", "application/json")

	if err != nil {
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(blogs)
}

func (handler *BlogHandler) GetAllComments(w http.ResponseWriter, r *http.Request) {
	comments, err := handler.BlogService.GetAllComments()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comments)
}
