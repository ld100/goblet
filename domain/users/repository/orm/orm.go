// GORM-based implementation of User Repository
package orm

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/lib/pq"

	models "github.com/ld100/goblet/domain/users"
	"github.com/ld100/goblet/domain/users/repository"
	"github.com/ld100/goblet/util/log"
)

type ormUserRepository struct {
	Conn *gorm.DB
}

func (repo *ormUserRepository) GetAll() ([]*models.User, error) {
	var users []*models.User
	var errs []error
	errs = repo.Conn.Find(&users).GetErrors()
	if len(errs) > 0 {
		log.Warn(errs)
		return nil, models.INTERNAL_SERVER_ERROR
	}
	return users, nil
}

func (repo *ormUserRepository) GetByID(id uint) (*models.User, error) {
	u := &models.User{ID: id}
	var errs []error
	errs = repo.Conn.First(&u, u.ID).GetErrors()
	if len(errs) > 0 {
		log.Debug(errs)
		return nil, models.NOT_FOUND_ERROR
	}

	return u, nil
}

func (repo *ormUserRepository) GetByEmail(email string) (*models.User, error) {
	u := &models.User{Email: email}
	var errs []error
	errs = repo.Conn.Where("email = ?", u.Email).First(&u).GetErrors()
	if len(errs) > 0 {
		log.Debug(errs)
		return nil, models.NOT_FOUND_ERROR
	}
	return u, nil

}

func (repo *ormUserRepository) Update(u *models.User) (*models.User, error) {
	var errs []error
	errs = repo.Conn.Save(&u).GetErrors()
	if len(errs) > 0 {
		log.Debug(errs)
		return nil, models.NOT_FOUND_ERROR
	}
	return u, nil
}

func (repo *ormUserRepository) Store(u *models.User) (uint, error) {
	if repo.Conn.NewRecord(u) {
		var errs []error
		errs = repo.Conn.Save(&u).GetErrors()
		if len(errs) > 0 {
			log.Debug(errs)
			// TODO: Check whether it is validation, database or data error
			return 0, models.CONFLICT_ERROR
		}
	}
	return u.ID, nil
}

func (repo *ormUserRepository) Delete(id uint) (bool, error) {
	u := &models.User{ID: id}
	if !repo.Conn.NewRecord(u) {
		var errs []error
		errs = repo.Conn.Delete(&u).GetErrors()
		if len(errs) > 0 {
			log.Debug(errs)
			return false, models.NOT_FOUND_ERROR
		}
	}
	return true, nil
}

// Initiation method
func NewOrmUserRepository(Conn *gorm.DB) repository.UserRepository {
	return &ormUserRepository{Conn}
}
