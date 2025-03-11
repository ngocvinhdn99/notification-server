package routes

import (
	"draft-notification/controllers"

	"github.com/labstack/echo/v4"
)

func UserDeliveryServerRoute(e *echo.Echo) {
	e.POST("/user-delivery-server", controllers.CreateUserDeliveryServer)
	e.GET("/user-delivery-server", controllers.GetAllUserDeliveryServers)
	e.GET("/user-delivery-server/:id", controllers.GetUserDeliveryServerDetail)
	e.PUT("/user-delivery-server/:id", controllers.UpdateUserDeliveryServer)
	e.PATCH("/user-delivery-server/:id/change-status", controllers.ChangeStatusUserDeliveryServer)
}
