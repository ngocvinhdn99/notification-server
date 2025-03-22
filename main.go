package main

import (
	"draft-notification/configs"
	"draft-notification/middlewares"
	"draft-notification/routes"
	"log"

	"github.com/labstack/echo/v4"
)

// Main function
func main() {
	e := echo.New()

	e.Use(middlewares.ValidateToken)

	configs.ConnectDB()
	routes.WebviewServerRoute(e)
	routes.UserDeliveryServerRoute(e)
	routes.ConnectionRoute(e)

	log.Println("🚀 Server đang chạy trên http://localhost:8080")
	e.Start(":8080")
}

// package main

// import (
// 	"draft-notification/grpc"
// 	"draft-notification/queue"
// )

// func main() {
// 	queue.RunQueue()
// 	grpc.RunGrpc()
// }
