package routes

import (
	"rest/controller"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Routes() *gin.Engine {
	router := gin.Default()

	v1 := router.Group("/api/v1")
	{
		albums := v1.Group("/albums")
		{
			albums.GET(":id", controller.GetAlbumByID)
			albums.GET("", controller.GetAlbums)
			albums.POST("", controller.PostAlbum)
			albums.PATCH(":id", controller.UpdateAlbum)
			albums.DELETE(":id", controller.DeleteAlbumByID)
		}
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}
