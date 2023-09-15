package handlers

import (
	"apis/internal/dto"
	"apis/internal/entity"
	"apis/internal/infra/database"
	entitypkg "apis/pkg/entity"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

type ProductHandler struct {
	ProductDB database.ProductInterface
}

func NewProductHandler(db database.ProductInterface) *ProductHandler {
	return &ProductHandler{
		ProductDB: db,
	}
}

// Create Product godoc
// @Summary Cria um novo produto
// @Description Cria um novo produto
// @Tags products
// @Accept json
// @Produce json
// @Param request body dto.CreateProductInput true "Dados do produto"
// @Success 201 {object} dto.CreateProductOutput "Produto criado com sucesso"
// @Failure 400 {object} Error "Dados inválidos"
// @Failure 500 {object} Error "Erro interno"
// @Router /products [post]
// @Security ApiKeyAuth
func (ph *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	input := dto.CreateProductInput{}
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(Error{Message: err.Error()})
		return
	}

	p, err := entity.NewProduct(input.Name, input.Description, input.Price)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(Error{Message: err.Error()})
		return
	}

	err = ph.ProductDB.Create(p)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(Error{Message: err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(dto.CreateProductOutput{
		ID:          p.ID.String(),
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(Error{Message: err.Error()})
		return
	}
}

// GetProduct godoc
// @Summary Busca um produto
// @Description Busca um produto
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "ID do produto" Format(uuid)
// @Success 200 {object} entity.Product "Produto encontrado com sucesso"
// @Failure 404 {object} Error "Produto não encontrado"
// @Failure 500 {object} Error "Erro interno"
// @Router /products/{id} [get]
// @Security ApiKeyAuth
func (ph *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if _, rr := entitypkg.ParseID(id); rr != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(Error{Message: "id is required"})
		return
	}

	product, err := ph.ProductDB.FindById(id)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(Error{Message: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(product)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(Error{Message: err.Error()})
		return
	}
}

// UpdateProduct godoc
// @Summary Atualiza um produto
// @Description Atualiza um produto
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "ID do produto" Format(uuid)
// @Param request body dto.UpdateProductInput true "Dados do produto"
// @Success 202 {object} dto.UpdateProductOutput "Produto atualizado com sucesso"
// @Failure 400 {object} Error "Dados inválidos"
// @Failure 404 {object} Error "Produto não encontrado"
// @Failure 500 {object} Error "Erro interno"
// @Router /products/{id} [put]
// @Security ApiKeyAuth
func (ph *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if _, rr := entitypkg.ParseID(id); rr != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(Error{Message: "id is required"})
		return
	}

	input := dto.UpdateProductInput{}
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(Error{Message: err.Error()})
		return
	}

	product, err := ph.ProductDB.FindById(id)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(Error{Message: err.Error()})
		return
	}

	product.Name = input.Name
	product.Description = input.Description
	product.Price = input.Price

	err = ph.ProductDB.Update(product)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(Error{Message: err.Error()})
		return
	}

	w.WriteHeader(http.StatusAccepted)
	err = json.NewEncoder(w).Encode(dto.UpdateProductOutput{
		ID:          product.ID.String(),
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(Error{Message: err.Error()})
		return
	}
}

// DeleteProduct godoc
// @Summary Deleta um produto
// @Description Deleta um produto
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "ID do produto" Format(uuid)
// @Success 204 "Produto deletado com sucesso"
// @Failure 400 {object} Error "Dados inválidos"
// @Failure 404 {object} Error "Produto não encontrado"
// @Failure 500 {object} Error "Erro interno"
// @Router /products/{id} [delete]
// @Security ApiKeyAuth
func (ph *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if _, rr := entitypkg.ParseID(id); rr != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(Error{Message: "id is required"})
		return
	}

	err := ph.ProductDB.Delete(id)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(Error{Message: err.Error()})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetProducts godoc
// @Summary Busca todos os produtos
// @Description Busca todos os produtos
// @Tags products
// @Accept json
// @Produce json
// @Param page query string true "Número da página" default(1)
// @Param limit query string true "Número de itens por página" default(10)
// @Param sort query string true "Ordenação" default(asc)
// @Success 200 {array} entity.Product "Produtos encontrados com sucesso"
// @Failure 404 {object} Error "Produtos não encontrados"
// @Failure 500 {object} Error "Erro interno"
// @Router /products [get]
// @Security ApiKeyAuth
func (ph *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	pageParam := chi.URLParam(r, "page")
	limitParam := chi.URLParam(r, "limit")
	sort := chi.URLParam(r, "sort")

	page, err := strconv.Atoi(pageParam)
	if err != nil {
		page = 0
	}

	limit, err := strconv.Atoi(limitParam)
	if err != nil {
		limit = 0
	}

	products, err := ph.ProductDB.FindAll(page, limit, sort)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(Error{Message: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(products)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(Error{Message: err.Error()})
		return
	}
}
