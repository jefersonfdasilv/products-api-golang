package handlers

import (
	"apis/internal/dto"
	"apis/internal/entity"
	"apis/internal/infra/database"
	entitypkg "apis/pkg/entity"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type UserHandler struct {
	UserDB database.UserInterface
}

func NewUserHandler(db database.UserInterface) *UserHandler {
	return &UserHandler{
		UserDB: db,
	}
}

func (uh *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	input := dto.CreteUserInput{}
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user, err := entity.NewUser(input.Name, input.Email, input.Password)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = uh.UserDB.Create(user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (uh *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := uh.UserDB.FindAll()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(users)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (uh *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if _, rr := entitypkg.ParseID(id); rr != nil {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	user, err := uh.UserDB.FindById(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(user)
}

func (uh *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if _, rr := entitypkg.ParseID(id); rr != nil {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	input := dto.UpdateUserInput{}
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := uh.UserDB.FindById(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	user.
		SetName(input.Name).
		SetEmail(input.Email)
	if input.Password != "" {
		err = user.SetPassword(input.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	err = uh.UserDB.Update(user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (uh *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if _, rr := entitypkg.ParseID(id); rr != nil {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	_, err := uh.UserDB.FindById(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	err = uh.UserDB.Delete(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
