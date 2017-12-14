package repository

import "github.com/ld100/goblet/domain/users/models"

type UserRepository interface {
	//Fetch(cursor string, num int64) ([]*models.User, error)
	GetAll() ([]*models.User, error)
	GetByID(id uint) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Update(u *models.User) (*models.User, error)
	Store(u *models.User) (uint, error)
	Delete(id uint) (bool, error)
}
