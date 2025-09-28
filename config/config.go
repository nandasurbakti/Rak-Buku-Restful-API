package config

import (
	"log"

	"github.com/spf13/viper"
)

func InitConfig() {
	// Set nilai default untuk development lokal
	viper.SetDefault("PGHOST", "localhost")
	viper.SetDefault("PGPORT", 5432)
	viper.SetDefault("PGUSER", "postgres")
	viper.SetDefault("PGPASSWORD", "admin123")
	viper.SetDefault("PGDATABASE", "postgres")
	viper.SetDefault("PGSSLMODE", "disable")
	viper.SetDefault("JWT_SECRET", "mysecret")

	viper.SetConfigFile(".env") 
	viper.AutomaticEnv()        

	if err := viper.ReadInConfig(); err != nil {
		
		log.Println("Error: .env file tidak ditemukan")
	}
}