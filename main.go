package main

import (
	"log"
	"net/http"

	articlehandler "article/articleshare/handlers/articlehandler"

	"github.com/gin-gonic/gin"
)

// NewArticleConnection instance

func main() {

	gEngine := gin.Default()

	addHTTPHandlers(gEngine)

	if runErr := gEngine.Run(":8080"); runErr != nil {
		log.Fatalf("Failed to run server: %v", runErr)
	}

}

func addHTTPHandlers(router *gin.Engine) {
	versionAPI := router.Group("api/v0")
	versionAPI.POST("/list", articlehandler.ListArticlesHandler)
	versionAPI.POST("/create", articlehandler.CreateArticlesHandler)

}

func apiGatewayVersionHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"version": 1.0,
	})
}
