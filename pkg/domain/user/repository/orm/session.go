// GORM-based implementation of Session Repository
package orm

import (
	usererrors "github.com/ld100/goblet/pkg/domain/user/error"
	"github.com/ld100/goblet/pkg/domain/user/model"
	"github.com/ld100/goblet/pkg/domain/user/repository"
	"github.com/ld100/goblet/pkg/server/env"
)

type ormSessionRepository struct {
	Env *env.Env
}

func (repo *ormSessionRepository) GetAllByUser(u *model.User) ([]*model.Session, error) {
	log := repo.Env.Logger
	conn, err := repo.Env.DB.ORMConnection()
	if err != nil {
		return nil, err
	}

	var sessions []*model.Session
	var errs []error

	//errs = conn.Model(&u).Association("Sessions").Find(&sessions).GetErrors()
	errs = conn.Model(&u).Related(&sessions).GetErrors()
	if len(errs) > 0 {
		log.Warn(errs)
		return nil, usererrors.INTERNAL_SERVER_ERROR
	}
	return sessions, nil
}

func (repo *ormSessionRepository) GetByID(id uint) (*model.Session, error) {
	log := repo.Env.Logger
	conn, err := repo.Env.DB.ORMConnection()
	if err != nil {
		return nil, err
	}

	s := &model.Session{ID: id}
	var errs []error
	errs = conn.First(&s, s.ID).GetErrors()
	if len(errs) > 0 {
		log.Debug(errs)
		return nil, usererrors.NOT_FOUND_ERROR
	}

	return s, nil
}

func (repo *ormSessionRepository) GetByUuid(uuid string) (*model.Session, error) {
	log := repo.Env.Logger
	conn, err := repo.Env.DB.ORMConnection()
	if err != nil {
		return nil, err
	}

	s := &model.Session{Uuid: uuid}
	var errs []error
	errs = conn.Where("uuid = ?", s.Uuid).First(&s).GetErrors()
	if len(errs) > 0 {
		log.Debug(errs)
		return nil, usererrors.NOT_FOUND_ERROR
	}
	return s, nil

}

func (repo *ormSessionRepository) Store(s *model.Session) (uint, error) {
	log := repo.Env.Logger
	conn, err := repo.Env.DB.ORMConnection()
	if err != nil {
		return 0, err
	}

	if conn.NewRecord(s) {
		var errs []error
		errs = conn.Save(&s).GetErrors()
		if len(errs) > 0 {
			log.Debug(errs)
			// TODO: Check whether it is validation, database or data error
			return 0, usererrors.CONFLICT_ERROR
		}
	}
	return s.ID, nil
}

func (repo *ormSessionRepository) Delete(id uint) (bool, error) {
	log := repo.Env.Logger
	conn, err := repo.Env.DB.ORMConnection()
	if err != nil {
		return false, err
	}

	s := &model.Session{ID: id}
	if !conn.NewRecord(s) {
		var errs []error
		errs = conn.Delete(&s).GetErrors()
		if len(errs) > 0 {
			log.Debug(errs)
			return false, usererrors.NOT_FOUND_ERROR
		}
	}
	return true, nil
}

// Initiation method
func NewOrmSessionRepository(Env *env.Env) repository.SessionRepository {
	return &ormSessionRepository{Env}
}
