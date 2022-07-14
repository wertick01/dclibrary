package processors

import (
	"context"
	"errors"

	"dclibrary.com/internals/app/db"
	"dclibrary.com/internals/app/models"
)

type UsersProcessor struct {
	storage *db.UsersStorage
}

func NewUsersProcessor(storage *db.UsersStorage) *UsersProcessor {
	processor := new(UsersProcessor)
	processor.storage = storage
	return processor
}

func (processor *UsersProcessor) CreateUser(user models.User) (int, error) {

	if user.Name == "" {
		return 0, errors.New("name should not be empty")
	}

	return processor.storage.CreateNewUser(user.Name, user.Phone, user.Mail, user.Hash, user.Role.RoleId)
}

func (processor *UsersProcessor) FindUser(id int) (models.User, error) {
	user, err := processor.storage.GetUserById(id)

	if user.UserId != id {
		return user, errors.New("user not found")
	}

	return user, nil

}

func (processor *UsersProcessor) ListUsers(ctx context.Context, nameFilter string) ([]models.User, error) {
	return processor.storage.GetUsersList(ctx, nameFilter), nil
}
