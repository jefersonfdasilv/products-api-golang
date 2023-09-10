package entity

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewProduct(t *testing.T) {
	type args struct {
		name        string
		description string
		price       float64
	}
	tests := []struct {
		name string
		args args
	}{
		{"Product test1", args{"product test1", "product description 1", 1}},
		{"Product test2", args{"product test2", "product description 2", 200}},
		{"Product test3", args{"product test3", "product description 3", 300.2}},
		{"Product test4", args{"product test4", "product description 4", 4003333}},
	}
	for _, tt := range tests {
		product, err := NewProduct(tt.args.name, tt.args.description, tt.args.price)
		assert.Nil(t, err, "error should be nil")
		assert.NotNil(t, product, "product should not be nil")
		assert.NotEmpty(t, product.ID, "ID should not be empty")
		assert.NotEmpty(t, product.Name, "name should not be empty")
		assert.NotEmpty(t, product.Description, "description should not be empty")
		assert.NotEmpty(t, product.Price, "price should not be empty")
		assert.Greater(t, product.Price, 0.0, "price should be greater than 0")
		assert.Equal(t, tt.args.name, product.Name, "name should be the same")
		assert.Equal(t, tt.args.description, product.Description, "description should be the same")
		assert.Equal(t, tt.args.price, product.Price, "price should be the same")
	}
}

func TestNewProduct_WhenInvalidPrice(t *testing.T) {
	type args struct {
		name        string
		description string
		price       float64
	}
	tests := []struct {
		name string
		args args
	}{
		{"Product test1", args{"product test1", "invalid product description 1", 0}},
		{"Product test2", args{"product test2", "invalid product description 2", -1}},
		{"Product test3", args{"product test3", "invalid product description 3", -33}},
		{"Product test4", args{"product test4", "invalid product description 4", 0.0}},
	}
	for _, tt := range tests {
		product, err := NewProduct(tt.args.name, tt.args.description, tt.args.price)
		assert.NotNil(t, err, "error should not be nil")
		assert.Nil(t, product, "product should be nil")
		assert.Equal(t, ErrInvalidPrice, err, "error should be the same")
	}
}
