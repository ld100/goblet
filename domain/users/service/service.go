package service

import (
	"time"

	"github.com/ld100/goblet/domain/users/models"
	"github.com/ld100/goblet/domain/users/repository"
	usererrors "github.com/ld100/goblet/domain/users/errors"
)

type UserService interface {
	//Fetch(cursor string, num uint) ([]*models.User, string, error)
	GetAll() ([]*models.User, error)
	GetByID(id uint) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Update(u *models.User) (*models.User, error)
	Store(*models.User) (*models.User, error)
	Delete(id uint) (bool, error)
}

// Implementation of UserService
type userService struct {
	userRepo repository.UserRepository
}

func (u *userService) GetAll() ([]*models.User, error) {
	return u.userRepo.GetAll()
}

func (u *userService) GetByID(id uint) (*models.User, error) {
	return u.userRepo.GetByID(id)
}

func (u *userService) Update(user *models.User) (*models.User, error) {
	user.UpdatedAt = time.Now()
	return u.userRepo.Update(user)
}

func (u *userService) GetByEmail(email string) (*models.User, error) {
	return u.userRepo.GetByEmail(email)
}

func (u *userService) Store(user *models.User) (*models.User, error) {
	// NOTE: this validation is excessive, since already done by model/ORM itself
	// Validation is placed here just for example purpose
	existingUser, _ := u.GetByEmail(user.Email)
	if existingUser != nil {
		return nil, usererrors.CONFLICT_ERROR
	}

	id, err := u.userRepo.Store(user)
	if err != nil {
		return nil, err
	}

	user.ID = id
	return user, nil
}

func (u *userService) Delete(id uint) (bool, error) {
	existingUser, _ := u.GetByID(id)

	if existingUser == nil {
		return false, usererrors.NOT_FOUND_ERROR
	}

	return u.userRepo.Delete(id)
}

// Initiation method
func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo}
}
