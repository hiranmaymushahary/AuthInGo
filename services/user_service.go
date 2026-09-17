package services

import (
	env "AuthInGo/config/env"
	db "AuthInGo/db/repositories"
	"AuthInGo/utils"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type UserService interface {
	GetUserById() error
	CreateUser() error
	LogInUser() (string, error)
}

type UserServiceImpl struct {
	userRepository db.UserRepository
}

// NewUserService creates and returns a new instance of UserService
func NewUserService(ur db.UserRepository) UserService {
	return &UserServiceImpl{
		userRepository: ur,
	}
}

func (u *UserServiceImpl) GetUserById() error {
	fmt.Println("Fetching user in UserService")
	u.userRepository.GetByID()
	return nil
}

func (u *UserServiceImpl) CreateUser() error {
	fmt.Println("Creating user in userService")
	password := "example_password"
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return err
	}
	u.userRepository.Create(
		"username_example_1",
		"user1@example.com",
		hashedPassword,
	)
	return nil
}

func (u *UserServiceImpl) LogInUser() (string, error) {
	email := "user1@example.com"
	password := "example_password"

	// Step 1: Fetch user by email
	user, err := u.userRepository.GetByEmail(email)
	if err != nil {
		fmt.Println("Error fetching user by email:", err)
		return "", err
	}

	// Step 2: Ensure user exists
	if user == nil {
		fmt.Println("No user found with the given email")
		return "", fmt.Errorf("no user found with email: %s", email)
	}

	// Step 3: Validate password
	isPasswordValid := utils.CheckPasswordHash(password, user.Password)
	if !isPasswordValid {
		fmt.Println("Password does not match")
		return "", fmt.Errorf("invalid password")
	}

	// Step 4: Generate JWT token upon successful authentication
	payload := jwt.MapClaims{
		"email": user.Email,
		"id":    user.Id,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	tokenString, err := token.SignedString([]byte(env.GetString("JWT_SECRET", "TOKEN")))
	if err != nil {
		fmt.Println("Error signing token:", err)
		return "", err
	}

	fmt.Println("JWT Token:", tokenString)
	return tokenString, nil
}