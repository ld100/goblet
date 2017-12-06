package repository

import models "github.com/ld100/goblet/domain/users"

type UserRepository interface {
	//Fetch(cursor string, num int64) ([]*models.User, error)
	GetAll() ([]*models.User, error)
	GetByID(id int64) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Update(u *models.User) (*models.User, error)
	Store(u *models.User) (int64, error)
	Delete(id int64) (bool, error)
}
