package repository

import (
	"a21hc3NpZ25tZW50/model"
	"errors"

	"gorm.io/gorm"
)

type UserRepository interface {
	GetUserByEmail(email string) (model.User, error)
	CreateUser(user model.User) (model.User, error)
	GetUserTaskCategory() ([]model.UserTaskCategory, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *userRepository {
	return &userRepository{db}
}

func (r *userRepository) GetUserByEmail(email string) (model.User, error) {
	var user model.User
	if result := r.db.Where("email = ?", email).First(&user); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return model.User{}, nil
		}

		return model.User{}, result.Error
	}

	return user, nil
}

func (r *userRepository) CreateUser(user model.User) (model.User, error) {
	err := r.db.Create(&user).Error
	if err != nil {
		return user, err
	}
	return user, nil
}

func (r *userRepository) GetUserTaskCategory() ([]model.UserTaskCategory, error) {
	var UserTaskCategories []model.UserTaskCategory

	result := r.db.Table("users").
		Select("users.id, users.fullname, users.email, tasks.title as task, tasks.deadline, tasks.priority, tasks.status, categories.name as category").
		Joins("left join tasks on tasks.user_id = users.id").
		Joins("left join categories on tasks.category_id = categories.id").
		Scan(&UserTaskCategories)
	if result.Error != nil {
		return nil, result.Error
	}
	return UserTaskCategories, nil
}
