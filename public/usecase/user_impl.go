package usecase

import (
	"errors"
	"log"

	"github.com/sandlayth/supplier-api/public/entity"
	"github.com/sandlayth/supplier-api/public/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserUseCaseImpl struct {
	repo repository.UserRepository
}

func NewUserUseCase(repo repository.UserRepository) *UserUseCaseImpl {
	return &UserUseCaseImpl{repo: repo}
}

func (uc *UserUseCaseImpl) Register(email, password, firstName, lastName string) (string, error) {
	// Check if the user with the given email already exists
	existingUser, err := uc.repo.FindByEmail(email)
	if err == nil && existingUser != nil {
		log.Print(err)
		return "", errors.New("user with the given email already exists")
	}

	// Hash the password before storing it
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	newUser := &entity.User{
		Email:     email,
		Password:  string(hashedPassword),
		FirstName: firstName,
		LastName:  lastName,
	}

	userID, err := uc.repo.Create(newUser)
	if err != nil {
		return "", err
	}

	return userID, nil
}

func (uc *UserUseCaseImpl) Login(email, password string) (string, error) {
	// Find user by email
	user, err := uc.repo.FindByEmail(email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	// Compare hashed password with the provided password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	// Generate and return an authentication token (@TODO later use JWT token)
	// For simplicity, let's return the user ID as the token
	return user.ID, nil
}

func (uc *UserUseCaseImpl) Logout(userID string) error {
	// Perform any necessary cleanup or token invalidation logic
	// For simplicity, no additional logic is added
	return nil
}

func (uc *UserUseCaseImpl) GetUserInfo(userID string) (*entity.User, error) {
	// Find user by ID
	user, err := uc.repo.FindById(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Omit sensitive information before returning user data
	user.Password = "" // Do not expose the hashed password

	return user, nil
}

func (uc *UserUseCaseImpl) ListAll() (*[]entity.User, error) {
	// Find user by ID
	users, err := uc.repo.ListAll()
	if err != nil {
		return nil, err
	}

	return users, nil
}
