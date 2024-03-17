package handler

import (
	"database-example/model"
	"database-example/service"
	"encoding/json"
	"net/http"
)

type UserHandler struct {
	UserService *service.UserService
}

// func (handler *StudentHandler) Get(writer http.ResponseWriter, req *http.Request) {
// 	id := mux.Vars(req)["id"]
// 	log.Printf("Student sa id-em %s", id)
// 	// student, err := handler.StudentService.FindStudent(id)
// 	// wr(iter.Header).Set("Content-Type", "application/json")
// 	// if err != nil {
// 	// 	writer.WriteHeader(http.StatusNotFound)
// 	// 	return
// 	// }
// 	writer.WriteHeader(http.StatusOK)
// 	// json.NewEncoder(writer).Encode(student)
// }

func (handler *UserHandler) Create(writer http.ResponseWriter, req *http.Request) {
	var user model.User
	err := json.NewDecoder(req.Body).Decode(&user)

	// bodyBytes, err := ioutil.ReadAll(req.Body)
	// fmt.Println(string(bodyBytes))

	if err != nil {
		println("Error while parsing json")
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	err = handler.UserService.Create(&user)
	if err != nil {
		println("Error while creating a new user")
		writer.WriteHeader(http.StatusExpectationFailed)
		return
	}
	writer.WriteHeader(http.StatusCreated)
	writer.Header().Set("Content-Type", "application/json")
}
