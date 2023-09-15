package main

import (
	"apis/configs"
	_ "apis/docs"
	"apis/internal/entity"
	"apis/internal/infra/database"
	"apis/internal/infra/webserver/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
)

// @title API de Produtos
// @description API de Produtos do curso de Go da FullCycle
// @version 1.0
// @termsOfService http://localhost:8000/terms
// @contactName Jeferson Silva
// @contactEmail jeferson.mab@gmail.com
// @licenseUrl http://localhost:8000/license
// @licenseName MIT
// @host localhost:8000
// @schemes http https
// @basePath /
// @securityDefinitions.apiKey ApiKeyAuth
// @in header
// @name Authorization
// @tokenUrl http://localhost:8000/auth/generate_token
func main() {
	// Carrega a configuração
	conf, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	// Conecta-se ao banco de dados SQLite
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// Executa migrações automáticas
	err = db.AutoMigrate(&entity.User{}, &entity.Product{})
	if err != nil {
		panic(err)
	}

	// Inicializa os handlers
	productHandler := handlers.NewProductHandler(database.NewProduct(db))
	userHandler := handlers.NewUserHandler(database.NewUser(db))

	// Cria um roteador Chi
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.WithValue("jwtAuth", conf.JwtAuth))
	r.Use(middleware.WithValue("jwtExpiresIn", conf.JWTExpiresIn))
	// Rotas para produtos
	r.Route("/products", func(r chi.Router) {
		r.Use(ZapLoggerMiddleware(zap.NewExample()))
		r.Use(jwtauth.Verifier(conf.JwtAuth))
		r.Use(jwtauth.Authenticator)
		r.Post("/", productHandler.Create)
		r.Get("/", productHandler.GetProducts)
		r.Get("/{id}", productHandler.GetProduct)
		r.Put("/{id}", productHandler.UpdateProduct)
		r.Delete("/{id}", productHandler.DeleteProduct)
	})

	// Rotas para usuários
	r.Route("/users", func(r chi.Router) {
		r.Use(jwtauth.Verifier(conf.JwtAuth))
		r.Use(jwtauth.Authenticator)
		r.Patch("/{id}", userHandler.UpdateUser)
		r.Delete("/{id}", userHandler.DeleteUser)
		r.Get("/", userHandler.GetUsers)
		r.Get("/{id}", userHandler.GetUser)
	})
	r.Post("/users", userHandler.Create)
	r.Post("/users/auth/generate_token", userHandler.GetJwt)
	r.Get("/docs/*", httpSwagger.Handler(httpSwagger.URL("http://localhost:8000/docs/doc.json")))
	// Inicia o servidor HTTP na porta 8000
	err = http.ListenAndServe(":8000", r)
	if err != nil {
		panic(err)
	}
}

func ZapLoggerMiddleware(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Info("request started",
				zap.String("method", r.Method),
				zap.String("url", r.URL.String()),
				zap.String("remote_addr", r.RemoteAddr),
			)
			next.ServeHTTP(w, r)
			logger.Info("request completed")
		})
	}
}
