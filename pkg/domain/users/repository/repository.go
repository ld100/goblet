package repository

import "github.com/ld100/goblet/pkg/domain/users/models"

type UserRepository interface {
	//Fetch(cursor string, num int64) ([]*models.User, error)
	GetAll() ([]*models.User, error)
	GetByID(id uint) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Update(u *models.User) (*models.User, error)
	Store(u *models.User) (uint, error)
	Delete(id uint) (bool, error)
}

type SessionRepository interface {
	GetAllByUser(u *models.User) ([]*models.Session, error)
	GetByID(id uint) (*models.Session, error)
	GetByUuid(uuid string) (*models.Session, error)
	Store(s *models.Session) (uint, error)
	Delete(id uint) (bool, error)
}