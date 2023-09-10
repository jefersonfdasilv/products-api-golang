package entity

import (
	"apis/pkg/entity"
	"time"
)

var (
	ErrNameIsRequired        = entity.NewError("name is required")
	ErrDescriptionIsRequired = entity.NewError("description is required")
	ErrInvalidPrice          = entity.NewError("invalid price")
)

type Product struct {
	ID          entity.ID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	CreatedAt   time.Time `json:"created_at"`
}

func NewProduct(name, description string, price float64) (*Product, error) {
	product := &Product{
		ID:          entity.NewID(),
		Name:        name,
		Description: description,
		Price:       price,
		CreatedAt:   time.Now(),
	}
	err := ValidateProduct(product)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func ValidateProduct(p *Product) error {
	if p == nil {
		return entity.ErrInvalidEntity
	}
	if _, err := entity.ParseID(p.ID.String()); err != nil {
		return entity.ErrInvalidID
	}
	if p.ID.String() == "" {
		return entity.ErrIDIsRequire
	}
	if p.Name == "" {
		return ErrNameIsRequired
	}
	if p.Description == "" {
		return ErrDescriptionIsRequired
	}
	if p.Price <= 0 {
		return ErrInvalidPrice
	}
	return nil
}
