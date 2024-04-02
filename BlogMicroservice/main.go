package main

import (
	"database-example/handler"
	"database-example/model"
	"database-example/repo"
	"database-example/service"

	"log"
	"net/http"

	"github.com/gorilla/handlers"
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

	database.AutoMigrate(&model.Student{})
	database.AutoMigrate(&model.Blog{})
	database.AutoMigrate(&model.Comment{})

	// ODKOMENTARISATI LINIJE "database.Exec" KOJE NEMAS SACUVANE U BAZI
	// da bi bile dodate, one koje vec imas, zakomentarisi jer ce pokusati opet da doda - primary key constraint

	// database.Exec("INSERT INTO students VALUES ('aec7e123-233d-4a09-a289-75308ea5b7e6', 'Marko Markovic', 'Graficki dizajn')")

	// publishDate := time.Date(2024, time.March, 18, 12, 0, 0, 0, time.Local)
	// query := fmt.Sprintf("insert into blogs values ('33686a82-6686-4d40-99b3-f0736c2bc7f4', "+ // id
	// 	"2, "+ // user id
	// 	"'Is xiaomi a good phone?', "+ // title
	// 	"'Everyone is buying a xiaomi phone nowadays, so Im wondering if they are actually worth buying and how long they last', "+ // description
	// 	"'%s', '%d')", publishDate.Format("2006-01-02 15:04:05"), model.BlogStatus(1)) // publish date, status
	// database.Exec(query)

	// publishDate2 := time.Date(2024, time.March, 19, 12, 0, 0, 0, time.Local)
	// query2 := fmt.Sprintf("insert into comments values ('2e998703-78dd-4076-8cf4-b8bb7e19e500', "+ // id
	// 	"3, "+ // user id
	// 	"'33686a82-6686-4d40-99b3-f0736c2bc7f4', "+ // blog id
	// 	"'%s', 'I personaly think that all chineese phones are trash', '%s')", // publish date, text, last change date
	// 	publishDate2.Format("2006-01-02 15:04:05"), publishDate2.Format("2006-01-02 15:04:05"))
	// database.Exec(query2)

	return database
}

func startServer( /*handler *handler.StudentHandler,*/ handler1 *handler.BlogHandler, handler2 *handler.BlogHandler) {
	router := mux.NewRouter().StrictSlash(true)

	// router.HandleFunc("/students/{id}", handler.Get).Methods("GET")
	// router.HandleFunc("/students", handler.Create).Methods("POST")
	router.HandleFunc("/blogs/create", handler1.CreateBlog).Methods("POST")
	router.HandleFunc("/comments/create", handler2.CreateComment).Methods("POST")
	router.HandleFunc("/feed", handler1.FindAllBlogs).Methods("GET")

	router.HandleFunc("/blogs/{blogId}/comments", handler2.GetAllCommentsByBlogId).Methods("GET")
	router.HandleFunc("/comments", handler2.GetAllComments).Methods("GET")

	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static")))
	println("Server starting")
	log.Fatal(http.ListenAndServe(":8081",
		handlers.CORS(
			handlers.AllowedOrigins([]string{"*"}),
			handlers.AllowedMethods([]string{"GET", "POST", "PUT", "OPTIONS"}),
			handlers.AllowedHeaders([]string{"Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With"}),
		)(router)))

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

	repo1 := &repo.BlogRepository{DatabaseConnection: database}
	service1 := &service.BlogService{BlogRepo: repo1}
	handler1 := &handler.BlogHandler{BlogService: service1}

	repo2 := &repo.CommentRepository{DatabaseConnection: database}
	service2 := &service.BlogService{CommentRepo: repo2}
	handler2 := &handler.BlogHandler{BlogService: service2}

	startServer( /*handler,*/ handler1, handler2)
}
