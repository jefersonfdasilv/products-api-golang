package dto

type CreateProductInput struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

type UpdateProductInput struct {
	CreateProductInput
}

type CreteUserInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateUserInput struct {
	CreteUserInput
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
