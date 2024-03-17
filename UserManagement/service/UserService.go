package service

import (
	"database-example/model"
	"database-example/repo"
)

type UserService struct {
	UserRepo *repo.UserRepository
}

// func (service *StudentService) FindUser(id string) (*model.Student, error) {
// 	student, err := service.StudentRepo.FindById(id)
// 	if err != nil {
// 		return nil, fmt.Errorf(fmt.Sprintf("menu item with id %s not found", id))
// 	}
// 	return &student, nil
// }

func (service *UserService) Create(user *model.User) error {
	err := service.UserRepo.CreateUser(user)
	if err != nil {
		return err
	}
	return nil
}
