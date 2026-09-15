package services

import (
	db "AuthInGo/db/repositories"
	"fmt"
)

type UserService interface {
	GetUserById() error
}

type UserServiceImpl struct {
	userRepository db.UserRepository
}

func (u *UserServiceImpl) GetUserById() error {
	fmt.Println("Fetching user in UserService")
	u.userRepository.GetByID()
	return nil
}

// NewUserService creates and returns a new instance of UserService
func NewUserService(ur db.UserRepository) UserService {
	return &UserServiceImpl{
		userRepository: ur,
	}
}
