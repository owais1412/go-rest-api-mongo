package routes

import (
	"rest/controller"

	"github.com/gin-gonic/gin"
)

func Routes() *gin.Engine {
	router := gin.Default()
	router.GET("/albums", controller.GetAlbums)
	router.GET("/albums/:id", controller.GetAlbumByID)
	router.POST("/albums", controller.PostAlbum)
	router.PATCH("/albums/:id", controller.UpdateAlbum)
	router.DELETE("/albums/:id", controller.DeleteAlbumByID)

	return router
}
