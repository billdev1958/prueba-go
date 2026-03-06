package main

import (
	"log"
	"prueba-go/internal/app"
)

// @title Commerce API
// @version 1.0
// @description API for commerce, transactions, reports and audit logs.
// @termsOfService http://swagger.io/terms/

// @host localhost:8080
// @BasePath /api/v1
// @query.collection.format multi

// @securityDefinitions.apikey UserID
// @in header
// @name X-User-Id

func main() {
	application, err := app.NewApp()
	if err != nil {
		log.Fatal(err)
	}

	if err := application.Run(); err != nil {
		log.Fatal(err)
	}
}
