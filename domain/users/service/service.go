package repository

import (
	"time"

	models "github.com/ld100/goblet/domain/users"
	"github.com/ld100/goblet/domain/users/repository"
)

type UserService interface {
	//Fetch(cursor string, num int64) ([]*models.User, string, error)
	GetAll() ([]*models.User, string, error)
	GetByID(id int64) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Update(u *models.User) (*models.User, error)
	Store(*models.User) (*models.User, error)
	Delete(id int64) (bool, error)
}

// Implementation of UserService
type userService struct {
	userRepo repository.UserRepository
}

func (u *userService) GetAll() ([]*models.User, error) {
	return u.userRepo.GetAll(id)
}

func (u *userService) GetByID(id int64) (*models.User, error) {
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
		return nil, models.CONFLIT_ERROR
	}

	id, err := u.userRepo.Store(user)
	if err != nil {
		return nil, err
	}

	user.ID = id
	return user, nil
}

func (u *userService) Delete(id int64) (bool, error) {
	existingUser, _ := u.GetByID(id)

	if existingUser == nil {
		return false, models.NOT_FOUND_ERROR
	}

	return u.userRepo.Delete(id)
}

// Initiation method
func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo}
}
