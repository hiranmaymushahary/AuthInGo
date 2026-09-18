package services

import (
    env "AuthInGo/config/env"
    db "AuthInGo/db/repositories"
    "AuthInGo/dto"
    "AuthInGo/models"
    "AuthInGo/utils"
    "fmt"

    "github.com/golang-jwt/jwt/v5"
)

type UserService interface {
    GetUserById(id string) (*models.User, error)
    CreateUser(payload *dto.CreateUserRequestDTO) (*models.User, error)
    LogInUser(payload *dto.LogInUserRequestDTO) (string, error)
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

func (u *UserServiceImpl) GetUserById(id string) (*models.User, error) {
    fmt.Println("Fetching user in UserService")

    // Call the repository layer to fetch the user by ID and capture the result
    user, err := u.userRepository.GetByID(id)

    if err != nil {
        fmt.Println("Error fetching user:", err)
        return nil, err
    }
    return user, nil
}

func (u *UserServiceImpl) CreateUser(payload *dto.CreateUserRequestDTO) (*models.User, error) {
    fmt.Println("Creating user in userService")

    // Step 1. HASH THE PASSWORD USING utils.HashPassword
    hashedPassword, err := utils.HashPassword(payload.Password)
    if err != nil {
        fmt.Println("Error hashing Password:", err)
        return nil, err
    }

    // Step 2. CALL THE REPOSITORY TO CREATE THE USER
    // Fix: Used short variable assignment `:=` instead of comparison operator `!=`
    user, err := u.userRepository.Create(payload.Username, payload.Email, hashedPassword)
    if err != nil {
        fmt.Println("Error while creating user:", err)
        return nil, err
    }

    // RETURN THE CREATED USER
    return user, nil
}

func (u *UserServiceImpl) LogInUser(payload *dto.LogInUserRequestDTO) (string, error) {
    email := payload.Email
    password := payload.Password

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
    jwtPayload := jwt.MapClaims{
        "email": user.Email,
        "id":    user.Id,
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtPayload)

    tokenString, err := token.SignedString([]byte(env.GetString("JWT_SECRET", "TOKEN")))
    if err != nil {
        fmt.Println("Error signing token:", err)
        return "", err
    }

    fmt.Println("JWT Token:", tokenString)
    return tokenString, nil
}