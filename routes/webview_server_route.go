package routes

import (
	"draft-notification/controllers"

	"github.com/labstack/echo/v4"
)

func WebviewServerRoute(e *echo.Echo) {
	e.POST("/webview-server", controllers.CreateWebviewServer)
	e.GET("/webview-server", controllers.GetAllWebviewServers)
	e.GET("/webview-server/:id", controllers.GetWebviewServerDetail)
	e.PUT("/webview-server/:id", controllers.UpdateWebviewServer)
	e.PATCH("/webview-server/:id/change-status", controllers.ChangeStatusWebviewServer)
}
