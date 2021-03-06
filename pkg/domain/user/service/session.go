package service

import (
	usererrors "github.com/ld100/goblet/pkg/domain/user/error"
	"github.com/ld100/goblet/pkg/domain/user/model"
	"github.com/ld100/goblet/pkg/domain/user/repository"
)

type SessionService interface {
	GetAllByUser(u *model.User) ([]*model.Session, error)
	GetByID(id uint) (*model.Session, error)
	GetByUuid(uuid string) (*model.Session, error)
	Store(*model.Session) (*model.Session, error)
	Delete(id uint) (bool, error)
}

// Implementation of SessionService
type sessionService struct {
	sessionRepo repository.SessionRepository
}

func (s *sessionService) GetAllByUser(u *model.User) ([]*model.Session, error) {
	return s.sessionRepo.GetAllByUser(u)
}

func (s *sessionService) GetByID(id uint) (*model.Session, error) {
	return s.sessionRepo.GetByID(id)
}

func (s *sessionService) GetByUuid(uuid string) (*model.Session, error) {
	return s.sessionRepo.GetByUuid(uuid)
}

func (s *sessionService) Store(session *model.Session) (*model.Session, error) {
	id, err := s.sessionRepo.Store(session)
	if err != nil {
		return nil, err
	}

	session.ID = id
	return session, nil
}

func (s *sessionService) Delete(id uint) (bool, error) {
	existingSession, _ := s.GetByID(id)

	if existingSession == nil {
		return false, usererrors.NOT_FOUND_ERROR
	}

	return s.sessionRepo.Delete(id)
}

// Initiation method
func NewSessionService(repo repository.SessionRepository) SessionService {
	return &sessionService{repo}
}
