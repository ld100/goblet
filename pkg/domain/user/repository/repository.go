package repository

import "github.com/ld100/goblet/pkg/domain/user/model"

type UserRepository interface {
	//Fetch(cursor string, num int64) ([]*model.User, error)
	GetAll() ([]*model.User, error)
	GetByID(id uint) (*model.User, error)
	GetByEmail(email string) (*model.User, error)
	Update(u *model.User) (*model.User, error)
	Store(u *model.User) (uint, error)
	Delete(id uint) (bool, error)
}

type SessionRepository interface {
	GetAllByUser(u *model.User) ([]*model.Session, error)
	GetByID(id uint) (*model.Session, error)
	GetByUuid(uuid string) (*model.Session, error)
	Store(s *model.Session) (uint, error)
	Delete(id uint) (bool, error)
}