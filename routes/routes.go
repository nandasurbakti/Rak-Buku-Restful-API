package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nandasurbakti/Rak-Buku-Restful-API/controllers"
	"github.com/nandasurbakti/Rak-Buku-Restful-API/middleware"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupServer() *gin.Engine {
	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api")
	{
		users := api.Group("/users")
		{
			users.POST("/login", controllers.Login)
            users.POST("/register", controllers.Register)

			users.Use(middleware.JWTAuthMiddleware())
		}

		categories := api.Group("/categories")
		{
			categories.GET("", controllers.GetCategories)
			categories.GET("/:id", controllers.GetCategoryByID)
            categories.GET("/:id/books", controllers.GetBooksByCategory)

			categories.Use(middleware.JWTAuthMiddleware())
            categories.POST("", controllers.PostCategory)
            categories.PUT("/:id", controllers.UpdateCategory)
            categories.DELETE("/:id", controllers.DeleteCategory)
		}

		books := api.Group("/books")
        {
            books.GET("", controllers.GetBooks)
			books.GET("/:id", controllers.GetBookByID)

			books.Use(middleware.JWTAuthMiddleware())
            books.POST("", controllers.PostBook)
            books.PUT("/:id", controllers.UpdateBook)
            books.DELETE("/:id", controllers.DeleteBook)
        }
	}

	return r
}