package main

import (
	"os"
	"github.com/SebasPu/Transacciones/cmd/server/handler"
	"github.com/SebasPu/Transacciones/internal/transacciones"
	"github.com/SebasPu/Transacciones/pkg/store"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/SebasPu/Transacciones/docs"
	"github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title MELI Bootcamp API
// @version 1.0
// @description This API Handle MELI Transactions.
// @termsOfService https://developers.mercadolibre.com.ar/es_ar/terminos-y-condiciones

// @contact.name API Support
// @contact.url https://developers.mercadolibre.com.ar/support

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
func main() {
	_ = godotenv.Load()
	db := store.New(store.FileType, "./transactions.json")
	repo := transacciones.NewRepository(db)

	service := transacciones.NewService(repo)

	t := handler.NewTransaccion(service)
	router := gin.Default()

	docs.SwaggerInfo.Host = os.Getenv("HOST")
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	//agrupacion de routers iguales
	transacciones := router.Group("/transacciones")
	{
		transacciones.GET("/getAll", t.GetAll())
		transacciones.GET("/lastId", t.LastId())
		transacciones.POST("/store", t.Store())
		transacciones.PUT("/:id", t.Update())
		transacciones.PATCH("/:id", t.UpdateCod())
		transacciones.DELETE("/:id", t.Delete())
	}

	router.Run()
}
