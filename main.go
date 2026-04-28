package main

import (
	"fmt"

	"github.com/essaherlandy/payment-midtrans/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {

	godotenv.Load(".env")
	app := fiber.New()

	// CORS manual (biar React bisa akses)
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
	}))

	app.Post("/create-transaction", services.CreateTransaction)
	app.Post("/webhook", services.Webhook)
	app.Post("/midtrans/callback", services.MidtransCallback)

	fmt.Println("Server running on :8080")
	app.Listen(":8080")
}
