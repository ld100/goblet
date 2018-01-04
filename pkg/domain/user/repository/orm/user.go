// GORM-based implementation of User Repository
package orm

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/lib/pq"

	usererrors "github.com/ld100/goblet/pkg/domain/user/error"
	"github.com/ld100/goblet/pkg/domain/user/model"
	"github.com/ld100/goblet/pkg/domain/user/repository"
	"github.com/ld100/goblet/pkg/util/log"
)

type ormUserRepository struct {
	Conn *gorm.DB
}

func (repo *ormUserRepository) GetAll() ([]*model.User, error) {
	var users []*model.User
	var errs []error
	errs = repo.Conn.Find(&users).GetErrors()
	if len(errs) > 0 {
		log.Warn(errs)
		return nil, usererrors.INTERNAL_SERVER_ERROR
	}
	return users, nil
}

func (repo *ormUserRepository) GetByID(id uint) (*model.User, error) {
	u := &model.User{ID: id}
	var errs []error
	errs = repo.Conn.First(&u, u.ID).GetErrors()
	if len(errs) > 0 {
		log.Debug(errs)
		return nil, usererrors.NOT_FOUND_ERROR
	}

	return u, nil
}

func (repo *ormUserRepository) GetByEmail(email string) (*model.User, error) {
	u := &model.User{Email: email}
	var errs []error
	errs = repo.Conn.Where("email = ?", u.Email).First(&u).GetErrors()
	if len(errs) > 0 {
		log.Debug(errs)
		return nil, usererrors.NOT_FOUND_ERROR
	}
	return u, nil

}

func (repo *ormUserRepository) Update(u *model.User) (*model.User, error) {
	var errs []error
	errs = repo.Conn.Save(&u).GetErrors()
	if len(errs) > 0 {
		log.Debug(errs)
		return nil, usererrors.NOT_FOUND_ERROR
	}
	return u, nil
}

func (repo *ormUserRepository) Store(u *model.User) (uint, error) {
	if repo.Conn.NewRecord(u) {
		var errs []error
		errs = repo.Conn.Save(&u).GetErrors()
		if len(errs) > 0 {
			log.Debug(errs)
			// TODO: Check whether it is validation, database or data error
			return 0, usererrors.CONFLICT_ERROR
		}
	}
	return u.ID, nil
}

func (repo *ormUserRepository) Delete(id uint) (bool, error) {
	u := &model.User{ID: id}
	if !repo.Conn.NewRecord(u) {
		var errs []error
		errs = repo.Conn.Delete(&u).GetErrors()
		if len(errs) > 0 {
			log.Debug(errs)
			return false, usererrors.NOT_FOUND_ERROR
		}
	}
	return true, nil
}

// Initiation method
func NewOrmUserRepository(Conn *gorm.DB) repository.UserRepository {
	return &ormUserRepository{Conn}
}
