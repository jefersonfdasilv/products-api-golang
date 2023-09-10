package database

import "apis/internal/entity"

type UserInterface interface {
	FindAll() ([]*entity.User, error)
	FindById(id string) (*entity.User, error)
	FindByEmail(email string) (*entity.User, error)
	Create(user *entity.User) error
	Update(user *entity.User) error
	Delete(id string) error
}

type ProductInterface interface {
	Create(product *entity.Product) error
	FindAll(page, limit int, sort string) ([]*entity.Product, error)
	FindById(id string) (*entity.Product, error)
	Update(product *entity.Product) error
	Delete(id string) error
}
