package main

import (
	"database-example/handler"
	"database-example/model"
	"database-example/repo"
	"database-example/service"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func initDB() *gorm.DB {
	dsn := "user=postgres password=ftn dbname=SOA host=localhost port=5432 sslmode=disable"
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		print(err)
		return nil
	}

	database.AutoMigrate(&model.Student{})
	database.AutoMigrate(&model.User{})
	database.Exec("INSERT IGNORE INTO students VALUES ('aec7e123-233d-4a09-a289-75308ea5b7e6', 'Marko Markovic', 'Graficki dizajn')")
	return database
}

func startServer(handler *handler.UserHandler) {
	router := mux.NewRouter().StrictSlash(true)

	// router.HandleFunc("/students/{id}", handler.Get).Methods("GET")
	router.HandleFunc("/users", handler.Create).Methods("POST")

	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static")))
	println("Server starting")
	log.Fatal(http.ListenAndServe(":8092", router))
}

func main() {
	database := initDB()
	if database == nil {
		print("FAILED TO CONNECT TO DB")
		return
	}
	// repo := &repo.StudentRepository{DatabaseConnection: database}
	// service := &service.StudentService{StudentRepo: repo}
	// handler := &handler.StudentHandler{StudentService: service}

	userRepo := &repo.UserRepository{DatabaseConnection: database}
	userService := &service.UserService{UserRepo: userRepo}
	handler := &handler.UserHandler{UserService: userService}

	startServer(handler)
}
