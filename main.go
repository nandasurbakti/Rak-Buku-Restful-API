package main

// @title Rak Buku API
// @version 1.0
// @description API RESTful untuk mengelola rak buku dengan kategori dan buku
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Ketik "Bearer" diikuti dengan spasi dan token JWT

import (
	"log"
	"os"

	_ "github.com/nandasurbakti/Rak-Buku-Restful-API/docs"

	_ "github.com/lib/pq"
	"github.com/nandasurbakti/Rak-Buku-Restful-API/config"
	"github.com/nandasurbakti/Rak-Buku-Restful-API/database"
	"github.com/nandasurbakti/Rak-Buku-Restful-API/database/connection"
	"github.com/nandasurbakti/Rak-Buku-Restful-API/routes"
)

func main() {
	config.InitConfig()

	// Koneksi database
	connection.ConnectDatabase()
	database.DBMigrate(connection.DB)

	// Setup routes
	r := routes.SetupServer()

	port := os.Getenv("PORT")
	if port == "" { 
		port = "8080"
	}

	log.Println("Server berjalan di port", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to run server:", err)
	}
}