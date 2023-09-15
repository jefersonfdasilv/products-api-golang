package configs

import (
	"fmt"
	"github.com/go-chi/jwtauth"
	"github.com/spf13/viper"
)

type Conf struct {
	DBDriver      string `mapstructure:"DB_DRIVER"`
	DbUser        string `mapstructure:"DB_USER"`
	DbPassword    string `mapstructure:"DB_PASSWORD"`
	DbHost        string `mapstructure:"DB_HOST"`
	DbPort        string `mapstructure:"DB_PORT"`
	DbName        string `mapstructure:"DB_NAME"`
	WebServerPort string `mapstructure:"WEB_SERVER_PORT"`
	JWTSecret     string `mapstructure:"JWT_SECRET"`
	JWTExpiresIn  int    `mapstructure:"JWT_EXPIRES_IN"`
	JwtAuth       *jwtauth.JWTAuth
}

func LoadConfig(path string) (*Conf, error) {
	var cfg *Conf
	viper.SetConfigName("app_config")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	cfg.JwtAuth = jwtauth.New("HS256", []byte(cfg.JWTSecret), nil)

	return cfg, nil
}
