package database

import (
	"apis/internal/entity"
	"gorm.io/gorm"
)

type User struct {
	DB *gorm.DB
}

func NewUser(db *gorm.DB) *User {
	return &User{DB: db}
}

func (u *User) Create(user *entity.User) error {
	return u.DB.Create(user).Error
}

func (u *User) FindAll() ([]*entity.User, error) {
	var users []*entity.User
	if err := u.DB.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (u *User) FindById(id string) (*entity.User, error) {
	var user entity.User
	if err := u.DB.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *User) FindByEmail(email string) (*entity.User, error) {
	var user entity.User
	if err := u.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *User) Update(user *entity.User) error {
	_, err := u.FindById(user.ID.String())
	if err != nil {
		return err
	}
	return u.DB.Save(user).Error
}

func (u *User) Delete(id string) error {
	user, err := u.FindById(id)
	if err != nil {
		return err
	}
	return u.DB.Delete(user).Error
}
