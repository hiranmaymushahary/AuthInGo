package services

import (
	db "AuthInGo/db/repositories"
	"fmt"
)

type UserService interface {
	Create() error
}

type UserServiceImpl struct {
	userRepository db.UserRepository
}

func (u *UserServiceImpl) Create() error {
	fmt.Println("Creating user in user service")
	u.userRepository.Create()
	return nil

}

func NewUserService(_userRepository db.UserRepository) UserService {
	return &UserServiceImpl{
		userRepository: _userRepository,
	}
}
