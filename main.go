package main

import (
	_ "rest/docs"
	"rest/middlewares"
	"rest/routes"
)

// @title REST API
// @version 1.0
// @description Album microservice server.
// @schemes http https

// @host      localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey  bearer
// @in                          header
// @name                        Authorization
func main() {

	r := routes.Routes()

	r.Run("localhost:" + middlewares.DotEnvVariable("PORT"))
}
