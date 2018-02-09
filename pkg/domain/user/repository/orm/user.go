// GORM-based implementation of User Repository
package orm

import (
	usererrors "github.com/ld100/goblet/pkg/domain/user/error"
	"github.com/ld100/goblet/pkg/domain/user/model"
	"github.com/ld100/goblet/pkg/domain/user/repository"
	"github.com/ld100/goblet/pkg/server/env"
)

type ormUserRepository struct {
	Env *env.Env
}

func (repo *ormUserRepository) GetAll() ([]*model.User, error) {
	log := repo.Env.Logger
	conn, err := repo.Env.DB.ORMConnection()
	if err != nil {
		return nil, err
	}

	var users []*model.User
	var errs []error
	errs = conn.Find(&users).GetErrors()
	if len(errs) > 0 {
		log.Warn(errs)
		return nil, usererrors.INTERNAL_SERVER_ERROR
	}
	return users, nil
}

func (repo *ormUserRepository) GetByID(id uint) (*model.User, error) {
	log := repo.Env.Logger
	conn, err := repo.Env.DB.ORMConnection()
	if err != nil {
		return nil, err
	}

	u := &model.User{ID: id}
	var errs []error
	errs = conn.First(&u, u.ID).GetErrors()
	if len(errs) > 0 {
		log.Debug(errs)
		return nil, usererrors.NOT_FOUND_ERROR
	}

	return u, nil
}

func (repo *ormUserRepository) GetByEmail(email string) (*model.User, error) {
	log := repo.Env.Logger
	conn, err := repo.Env.DB.ORMConnection()
	if err != nil {
		return nil, err
	}

	u := &model.User{Email: email}
	var errs []error
	errs = conn.Where("email = ?", u.Email).First(&u).GetErrors()
	if len(errs) > 0 {
		log.Debug(errs)
		return nil, usererrors.NOT_FOUND_ERROR
	}
	return u, nil

}

func (repo *ormUserRepository) Update(u *model.User) (*model.User, error) {
	log := repo.Env.Logger
	conn, err := repo.Env.DB.ORMConnection()
	if err != nil {
		return nil, err
	}

	var errs []error
	errs = conn.Save(&u).GetErrors()
	if len(errs) > 0 {
		log.Debug(errs)
		return nil, usererrors.NOT_FOUND_ERROR
	}
	return u, nil
}

func (repo *ormUserRepository) Store(u *model.User) (uint, error) {
	log := repo.Env.Logger
	conn, err := repo.Env.DB.ORMConnection()
	if err != nil {
		return 0, err
	}

	if conn.NewRecord(u) {
		var errs []error
		errs = conn.Save(&u).GetErrors()
		if len(errs) > 0 {
			log.Debug(errs)
			// TODO: Check whether it is validation, database or data error
			return 0, usererrors.CONFLICT_ERROR
		}
	}
	return u.ID, nil
}

func (repo *ormUserRepository) Delete(id uint) (bool, error) {
	log := repo.Env.Logger
	conn, err := repo.Env.DB.ORMConnection()
	if err != nil {
		return false, err
	}

	u := &model.User{ID: id}
	if !conn.NewRecord(u) {
		var errs []error
		errs = conn.Delete(&u).GetErrors()
		if len(errs) > 0 {
			log.Debug(errs)
			return false, usererrors.NOT_FOUND_ERROR
		}
	}
	return true, nil
}

// Initiation method
func NewOrmUserRepository(Env *env.Env) repository.UserRepository {
	return &ormUserRepository{Env}
}
