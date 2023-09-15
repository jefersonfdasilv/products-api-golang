package handlers

import (
	"apis/internal/dto"
	"apis/internal/entity"
	"apis/internal/infra/database"
	entitypkg "apis/pkg/entity"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth"
	"net/http"
	"time"
)

type UserHandler struct {
	UserDB       database.UserInterface
	Jwt          *jwtauth.JWTAuth
	JwtExpiresIn int
}

func NewUserHandler(db database.UserInterface) *UserHandler {
	return &UserHandler{
		UserDB: db,
	}
}

func (uh *UserHandler) GetAuth(r *http.Request) (*jwtauth.JWTAuth, int) {
	jwt := r.Context().Value("jwtAuth").(*jwtauth.JWTAuth)
	jwtExpiresIn := r.Context().Value("jwtExpiresIn").(int)
	return jwt, jwtExpiresIn
}

// GetJwt godoc
// @Summary Gera um token JWT
// @Description Gera um token JWT
// @Tags users
// @Accept json
// @Produce json
// @Param request body dto.LoginInput true "Credenciais dos usuário"
// @Success 200 {object} dto.GetJWTOutput "Usuário autenticado com sucesso"
// @Failure 400 {object} Error "Dados inválidos"
// @Failure 401 {object} Error "Credenciais inválidas"
// @Failure 500 {object} Error "Erro interno"
// @Router /users/auth/generate_token [post]
func (uh *UserHandler) GetJwt(w http.ResponseWriter, r *http.Request) {

	jwt, jwtExpiresIn := uh.GetAuth(r)

	input := dto.LoginInput{}
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(Error{Message: err.Error()})
		return
	}

	user, err := uh.UserDB.FindByEmail(input.Email)

	if err != nil {
		http.Error(w, "invalid credentials ", http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(Error{Message: err.Error()})
		return
	}

	if !user.CheckPassword(input.Password) {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(Error{Message: err.Error()})
		return
	}

	_, tokenString, _ := jwt.Encode(map[string]interface{}{
		"sub":   user.ID,
		"name":  user.Name,
		"email": user.Email,
		"exp":   time.Now().Add(time.Second * time.Duration(jwtExpiresIn)).Unix(),
	})

	accessToken := dto.GetJWTOutput{AccessToken: tokenString}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(accessToken)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(Error{Message: err.Error()})
	}
}

// Create user godoc
// @Summary Cria um usuário
// @Description Cria um usuário
// @Tags users
// @Accept json
// @Produce json
// @Param request body dto.CreteUserInput true "Dados do usuário"
// @Success 201 {object} dto.CreateUserOutput "Usuário criado com sucesso"
// @Failure 400 {object} Error "Dados inválidos"
// @Failure 500 {object} Error "Erro interno"
// @Router /users [post]
func (uh *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	input := dto.CreteUserInput{}
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(Error{Message: err.Error()})
		return
	}
	user, err := entity.NewUser(input.Name, input.Email, input.Password)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(Error{Message: err.Error()})
		return
	}

	err = uh.UserDB.Create(user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(Error{Message: err.Error()})
		return
	}
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(dto.CreateUserOutput{
		ID:    user.ID.String(),
		Name:  user.Name,
		Email: user.Email,
	})
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
