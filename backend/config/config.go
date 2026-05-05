package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Port              string `mapstructure:"PORT"`
	DBUsername        string `mapstructure:"DB_USERNAME"`
	DBPassword        string `mapstructure:"DB_PASSWORD"`
	DBUrl             string `mapstructure:"DB_URL"`
	RedisUrl          string `mapstructure:"REDIS_URL"`
	RedisPort         string `mapstructure:"REDIS_PORT"`
	DBDatabase        string `mapstructure:"DB_DATABASE"`
	DBPort            string `mapstructure:"DB_PORT"`
	MinioHost         string `mapstructure:"MINIO_HOST"`
	MinioPort         string `mapstructure:"MINIO_PORT"`
	MinioUser         string `mapstructure:"MINIO_USER"`
	MinioPassword     string `mapstructure:"MINIO_PASSWORD"`
	MinioBucket       string `mapstructure:"MINIO_BUCKET"`
	SecretKey         string `mapstructure:"SECRET_KEY"`
	FrontendURL       string `mapstructure:"FRONTEND_URL"`
	ElasticDBUser     string `mapstructure:"ELASTICDB_USER"`
	ElasticDBPassword string `mapstructure:"ELASTICDB_PASSWORD"`
	ElasticDBHost     string `mapstructure:"ELASTICDB_HOST"`
	ElasticDBPort     string `mapstructure:"ELASTICDB_PORT"`
	RajaOngkirAPIKey  string `mapstructure:"RAJAONGKIR_API_KEY"`
	RajaOngkirURL     string `mapstructure:"RAJAONGKIR_URL"`
}

var ENV Config

func LoadConfig() {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	viper.ReadInConfig()

	// viper.BindEnv("PORT")
	// viper.BindEnv("DB_USERNAME")
	// viper.BindEnv("DB_PASSWORD")
	// viper.BindEnv("DB_URL")
	// viper.BindEnv("REDIS_URL")
	// viper.BindEnv("REDIS_PORT")
	// viper.BindEnv("DB_DATABASE")
	// viper.BindEnv("DB_PORT")
	// viper.BindEnv("MINIO_HOST")
	// viper.BindEnv("MINIO_PORT")
	// viper.BindEnv("MINIO_USER")
	// viper.BindEnv("MINIO_PASSWORD")
	// viper.BindEnv("MINIO_BUCKET")
	// viper.BindEnv("SECRET_KEY")
	// viper.BindEnv("FRONTEND_URL")
	// viper.BindEnv("ELASTICDB_USER")
	// viper.BindEnv("ELASTICDB_PASSWORD")
	// viper.BindEnv("ELASTICDB_HOST")
	// viper.BindEnv("ELASTICDB_PORT")

	if err := viper.Unmarshal(&ENV); err != nil {
		log.Fatal("Cannot unmarshal config:", err)
	}
}
