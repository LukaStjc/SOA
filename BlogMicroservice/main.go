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
	//connectionStr := "root:root@tcp(localhost:3306)/students?charset=utf8mb4&parseTime=True&loc=Local"
	connectionParams := "user=postgres password=ftn dbname=SOA host=localhost port=5432 sslmode=disable"
	database, err := gorm.Open(postgres.Open(connectionParams), &gorm.Config{})

	if err != nil {
		print(err)
		return nil
	}

	if err := database.AutoMigrate(&model.Student{}); err != nil {
		log.Printf("AutoMigrate failed: %v", err)
		return nil
	}

	// ZAKOMENTARISATI SVE LINIJE "database.Exec" NAKON PRVOG POKRETANJA PROJEKTA
	// jer ce prvi put dodati red u tabelu a drugi put ce pokusati opet isto da doda pa ce biti Primary key constraint
	database.Exec("INSERT INTO students VALUES ('aec7e123-233d-4a09-a289-75308ea5b7e6', 'Marko Markovic', 'Graficki dizajn')")

	return database
}

func startServer(handler *handler.StudentHandler) {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/students/{id}", handler.Get).Methods("GET")
	router.HandleFunc("/students", handler.Create).Methods("POST")

	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static")))
	println("Server starting")
	log.Fatal(http.ListenAndServe(":8081", router))
}

func main() {
	database := initDB()
	if database == nil {
		print("FAILED TO CONNECT TO DB")
		return
	}
	repo := &repo.StudentRepository{DatabaseConnection: database}
	service := &service.StudentService{StudentRepo: repo}
	handler := &handler.StudentHandler{StudentService: service}

	startServer(handler)
}
