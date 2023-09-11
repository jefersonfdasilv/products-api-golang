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

func (ph *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	input := dto.CreateProductInput{}
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	p, err := entity.NewProduct(input.Name, input.Description, input.Price)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = ph.ProductDB.Create(p)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (ph *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if _, rr := entitypkg.ParseID(id); rr != nil {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	product, err := ph.ProductDB.FindById(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(product)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (ph *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if _, rr := entitypkg.ParseID(id); rr != nil {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	input := dto.UpdateProductInput{}
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	product, err := ph.ProductDB.FindById(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	product.Name = input.Name
	product.Description = input.Description
	product.Price = input.Price

	err = ph.ProductDB.Update(product)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (ph *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if _, rr := entitypkg.ParseID(id); rr != nil {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	err := ph.ProductDB.Delete(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

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
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(products)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
