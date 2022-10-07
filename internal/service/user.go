package service

import (
	"forum/internal/storage"
	"forum/models"
)

type User interface {
	GetUserByUsername(id string) (models.User, error)
}

type UserService struct {
	storage storage.User
}

func newUserService(storage storage.User) *UserService {
	return &UserService{
		storage: storage,
	}
}

func (u *UserService) GetUserByUsername(id string) (models.User, error) {
	var user models.User
	user, err := u.storage.GetUserByUsername(id)
	if err != nil {
		return user, err
	}
	return user, nil
}
