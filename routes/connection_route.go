package routes

import (
	"draft-notification/controllers"

	"github.com/labstack/echo/v4"
)

func ConnectionRoute(e *echo.Echo) {
	e.POST("/user-delivery-server/:userDeliveryServerId/connection", controllers.CreateConnection)
	e.GET("/user-delivery-server/:userDeliveryServerId/connections", controllers.GetAllConnections)
	e.PATCH("/connections/:id/update-web-hook-url", controllers.UpdateConnectionWebhookUrl)
	e.PATCH("/connections/:id/change-status", controllers.ChangeStatusConnection)
}
