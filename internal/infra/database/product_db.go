package database

import (
	"apis/internal/entity"
	"gorm.io/gorm"
)

type Product struct {
	DB *gorm.DB
}

func NewProduct(db *gorm.DB) *Product {
	return &Product{DB: db}
}

func (p *Product) Create(product *entity.Product) error {
	return p.DB.Create(product).Error
}

func (p *Product) FindById(id string) (*entity.Product, error) {
	var product entity.Product
	if err := p.DB.Where("id = ?", id).First(&product).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (p *Product) FindAll(page, limit int, sort string) ([]*entity.Product, error) {
	var products []*entity.Product

	if sort != "asc" && sort != "desc" {
		sort = "asc"
	}

	if page != 0 && limit != 0 {
		if err := p.DB.Order("created_at " + sort).Offset((page - 1) * limit).Limit(limit).Find(&products).Error; err != nil {
			return nil, err
		}
		return products, nil
	}

	if err := p.DB.Order("created_at " + sort).Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}

func (p *Product) Update(product *entity.Product) error {
	_, err := p.FindById(product.ID.String())
	if err != nil {
		return err
	}
	return p.DB.Save(product).Error
}

func (p *Product) Delete(id string) error {
	product, err := p.FindById(id)
	if err != nil {
		return err
	}
	return p.DB.Delete(product).Error
}
