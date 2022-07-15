package processors

import (
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

	return processor.storage.CreateNewUser(user.Name, user.Phone, user.Mail, user.Hash, user.RoleId)
}

func (processor *UsersProcessor) FindUser(id int64) (*models.User, error) {
	user, err := processor.storage.GetUserById(id)

	if err != nil {
		return user, errors.New("user not found")
	}

	return user, nil

}

func (processor *UsersProcessor) ListUsers() ([]*models.User, error) {
	return processor.storage.GetUsersList()
}

func (processor *UsersProcessor) UpdateUser(id int64) (*models.User, error) { //!!! ПРОВЕРИТЬ
	user, _ := processor.FindUser(id)

	changeduser, err := processor.storage.ChangeUserById(id, user)
	if err != nil {
		return user, errors.New("SOMETHING IS WRONG")
	}

	return changeduser, nil
}

func (processor *UsersProcessor) DeleteUser(id int64) (*models.User, error) {
	user, _ := processor.FindUser(id)
	_, err := processor.storage.DeleteUserById(id)
	if err != nil {
		return user, errors.New("CANNOT DELETE USER")
	}
	return user, nil
}
