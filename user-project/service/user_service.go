package service

import (
	"github.com/cagnosolutions/go-web-ddd/user-project/domain"
)

type UserService interface {
	GetUser(id int) (*domain.User, error)
	GetAllUsers() ([]*domain.User, error)
}

func NewUserService(repo domain.UserRepository) *DefaultUserService {
	return &DefaultUserService{repo: repo}
}

// DefaultUserService is the default implementation of the user service
type DefaultUserService struct {
	repo domain.UserRepository
}

func (s *DefaultUserService) GetUser(id int) (*domain.User, error) {
	return s.repo.Get(id)
}

func (s *DefaultUserService) GetAllUsers() ([]*domain.User, error) {
	return s.repo.GetAll()
}
