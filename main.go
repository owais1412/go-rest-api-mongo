package main

import (
	"rest/middlewares"
	"rest/routes"
)

func main() {

	router := routes.Routes()

	router.Run("localhost:" + middlewares.DotEnvVariable("PORT"))
}
