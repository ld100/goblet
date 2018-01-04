// GORM-based implementation of Session Repository
package orm

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/lib/pq"

	usererrors "github.com/ld100/goblet/pkg/domain/users/errors"
	"github.com/ld100/goblet/pkg/domain/users/models"
	"github.com/ld100/goblet/pkg/domain/users/repository"
	"github.com/ld100/goblet/pkg/util/log"
)

type ormSessionRepository struct {
	Conn *gorm.DB
}

func (repo *ormSessionRepository) GetAllByUser(u *models.User) ([]*models.Session, error) {
	var sessions []*models.Session
	var errs []error

	//errs = repo.Conn.Model(&u).Association("Sessions").Find(&sessions).GetErrors()
	errs = repo.Conn.Model(&u).Related(&sessions).GetErrors()
	if len(errs) > 0 {
		log.Warn(errs)
		return nil, usererrors.INTERNAL_SERVER_ERROR
	}
	return sessions, nil
}

func (repo *ormSessionRepository) GetByID(id uint) (*models.Session, error) {
	s := &models.Session{ID: id}
	var errs []error
	errs = repo.Conn.First(&s, s.ID).GetErrors()
	if len(errs) > 0 {
		log.Debug(errs)
		return nil, usererrors.NOT_FOUND_ERROR
	}

	return s, nil
}

func (repo *ormSessionRepository) GetByUuid(uuid string) (*models.Session, error) {
	s := &models.Session{Uuid: uuid}
	var errs []error
	errs = repo.Conn.Where("uuid = ?", s.Uuid).First(&s).GetErrors()
	if len(errs) > 0 {
		log.Debug(errs)
		return nil, usererrors.NOT_FOUND_ERROR
	}
	return s, nil

}

func (repo *ormSessionRepository) Store(s *models.Session) (uint, error) {
	if repo.Conn.NewRecord(s) {
		var errs []error
		errs = repo.Conn.Save(&s).GetErrors()
		if len(errs) > 0 {
			log.Debug(errs)
			// TODO: Check whether it is validation, database or data error
			return 0, usererrors.CONFLICT_ERROR
		}
	}
	return s.ID, nil
}

func (repo *ormSessionRepository) Delete(id uint) (bool, error) {
	s := &models.Session{ID: id}
	if !repo.Conn.NewRecord(s) {
		var errs []error
		errs = repo.Conn.Delete(&s).GetErrors()
		if len(errs) > 0 {
			log.Debug(errs)
			return false, usererrors.NOT_FOUND_ERROR
		}
	}
	return true, nil
}

// Initiation method
func NewOrmSessionRepository(Conn *gorm.DB) repository.SessionRepository {
	return &ormSessionRepository{Conn}
}
