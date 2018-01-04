package service

import (
	"time"

	usererrors "github.com/ld100/goblet/pkg/domain/user/error"
	"github.com/ld100/goblet/pkg/domain/user/model"
	"github.com/ld100/goblet/pkg/domain/user/repository"
)

type UserService interface {
	//Fetch(cursor string, num uint) ([]*model.User, string, error)
	GetAll() ([]*model.User, error)
	GetByID(id uint) (*model.User, error)
	GetByEmail(email string) (*model.User, error)
	Update(u *model.User) (*model.User, error)
	Store(*model.User) (*model.User, error)
	Delete(id uint) (bool, error)
}

// Implementation of SessionService
type userService struct {
	userRepo repository.UserRepository
}

func (u *userService) GetAll() ([]*model.User, error) {
	return u.userRepo.GetAll()
}

func (u *userService) GetByID(id uint) (*model.User, error) {
	return u.userRepo.GetByID(id)
}

func (u *userService) Update(user *model.User) (*model.User, error) {
	user.UpdatedAt = time.Now()
	return u.userRepo.Update(user)
}

func (u *userService) GetByEmail(email string) (*model.User, error) {
	return u.userRepo.GetByEmail(email)
}

func (u *userService) Store(user *model.User) (*model.User, error) {
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
