package main

import (
	"shuttle/databases"
	"shuttle/routes"
	"shuttle/utils"
	zerolog "shuttle/logger"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/spf13/viper"
)

func main() {
	utils.InitFirebase()
	zerolog.InitLogger()

	app := fiber.New()

	app.Use(cors.New())

	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${method} ${path} [${status}] ${latency}\n",
	}))

	db, err := databases.PostgresConnection()
	if err != nil {
		panic(err)
	}

	routes.Route(app, db)

	if err := app.Listen(viper.GetString("BASE_URL")); err != nil {
        panic(err)
    }
}