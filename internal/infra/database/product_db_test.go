package database

import (
	"apis/internal/entity"
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

func TestProduct_Create(t *testing.T) {
	db, cleanup := initDBTest(t, &entity.Product{})
	defer cleanup()

	product, err := entity.NewProduct("product test", "product description", 21.9)
	if err != nil {
		t.Error(err)
	}
	productDb := NewProduct(db)
	err = productDb.Create(product)
	assert.Nil(t, err, "error should be nil")
	assert.NotNil(t, product, "product should not be nil")
	assert.NotEmpty(t, product.ID, "product ID should not be empty")

	var productCreated entity.Product
	err = db.Where("id = ?", product.ID).First(&productCreated).Error
	assert.Nil(t, err, "error should be nil")
	assert.NotNil(t, productCreated, "product should not be nil")
	assert.Equal(t, product.ID, productCreated.ID, "ID should be the same")
	assert.Equal(t, product.Name, productCreated.Name, "name should be the same")
	assert.Equal(t, product.Description, productCreated.Description, "description should be the same")
	assert.Equal(t, product.Price, productCreated.Price, "price should be the same")
}

func TestProduct_FindAll(t *testing.T) {
	db, cleanup := initDBTest(t, &entity.Product{})
	defer cleanup()

	for i := 0; i < 106; i++ {
		product, err := entity.NewProduct(fmt.Sprintf("product test %d", i), fmt.Sprintf("product description %d", i), rand.Float64()*100.0)
		if err != nil {
			t.Error(err)
		}
		productDb := NewProduct(db)
		err = productDb.Create(product)
		assert.Nil(t, err, "error should be nil")
		assert.NotNil(t, product, "product should not be nil")
	}

	products, err := NewProduct(db).FindAll(1, 10, "asc")
	assert.Nil(t, err, "error should be nil")
	assert.NotNil(t, products, "products should not be nil")
	assert.Equal(t, 10, len(products), "length should be 10")
	assert.Equal(t, "product test 0", products[0].Name, "name should be the same")
	assert.Equal(t, "product test 9", products[9].Name, "name should be the same")

	products, err = NewProduct(db).FindAll(2, 10, "asc")
	assert.Nil(t, err, "error should be nil")
	assert.NotNil(t, products, "products should not be nil")
	assert.Equal(t, 10, len(products), "length should be 10")
	assert.Equal(t, "product test 10", products[0].Name, "name should be the same")
	assert.Equal(t, "product test 19", products[9].Name, "name should be the same")

	products, err = NewProduct(db).FindAll(11, 10, "asc")
	assert.Nil(t, err, "error should be nil")
	assert.NotNil(t, products, "products should not be nil")
	assert.Equal(t, 6, len(products), "length should be 10")
	assert.Equal(t, "product test 100", products[0].Name, "name should be the same")
	assert.Equal(t, "product test 105", products[5].Name, "name should be the same")

}

func TestProduct_FindById(t *testing.T) {
	db, cleanup := initDBTest(t, &entity.Product{})
	defer cleanup()

	product, err := entity.NewProduct("product test", "product description", 21.9)
	if err != nil {
		t.Error(err)
	}
	productDb := NewProduct(db)
	err = productDb.Create(product)
	assert.Nil(t, err, "error should be nil")
	assert.NotNil(t, product, "product should not be nil")
	assert.NotEmpty(t, product.ID, "product ID should not be empty")

	productCreated, err := productDb.FindById(product.ID.String())
	assert.Nil(t, err, "error should be nil")
	assert.NotNil(t, productCreated, "product should not be nil")
	assert.Equal(t, product.ID, productCreated.ID, "ID should be the same")
	assert.Equal(t, product.Name, productCreated.Name, "name should be the same")
	assert.Equal(t, product.Description, productCreated.Description, "description should be the same")
	assert.Equal(t, product.Price, productCreated.Price, "price should be the same")
}

func TestProduct_Update(t *testing.T) {
	db, cleanup := initDBTest(t, &entity.Product{})
	defer cleanup()

	product, err := entity.NewProduct("product test", "product description", 21.9)
	if err != nil {
		t.Error(err)
	}
	productDb := NewProduct(db)
	err = productDb.Create(product)
	assert.Nil(t, err, "error should be nil")
	assert.NotNil(t, product, "product should not be nil")
	assert.NotEmpty(t, product.ID, "product ID should not be empty")

	product.Name = "product test updated"
	product.Description = "product description updated"
	product.Price = 99.99
	err = productDb.Update(product)
	assert.Nil(t, err, "error should be nil")

	productCreated, err := productDb.FindById(product.ID.String())
	assert.Nil(t, err, "error should be nil")
	assert.NotNil(t, productCreated, "product should not be nil")
	assert.Equal(t, product.ID, productCreated.ID, "ID should be the same")
	assert.Equal(t, product.Name, productCreated.Name, "name should be the same")
	assert.Equal(t, product.Description, productCreated.Description, "description should be the same")
	assert.Equal(t, product.Price, productCreated.Price, "price should be the same")
}

func TestProduct_Delete(t *testing.T) {
	db, cleanup := initDBTest(t, &entity.Product{})
	defer cleanup()

	product, err := entity.NewProduct("product test", "product description", 21.9)
	if err != nil {
		t.Error(err)
	}
	productDb := NewProduct(db)
	err = productDb.Create(product)
	assert.Nil(t, err, "error should be nil")
	assert.NotNil(t, product, "product should not be nil")
	assert.NotEmpty(t, product.ID, "product ID should not be empty")

	err = productDb.Delete(product.ID.String())
	assert.Nil(t, err, "error should be nil")

	_, err = productDb.FindById(product.ID.String())
	assert.NotNil(t, err, "error should not be nil")
}
