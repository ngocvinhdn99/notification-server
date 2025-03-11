package helpers

import (
	"draft-notification/responses"
	"net/http"

	"github.com/labstack/echo/v4"
)

// Helper function to handle error response
func HandleError(c echo.Context, status int, errMsg string) error {
	return c.JSON(status, responses.Response{Code: status, Message: errMsg, Data: &echo.Map{}})
}

// Helper function for success response
func HandleSuccess(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusOK, responses.Response{Code: http.StatusOK, Message: "success", Data: data})
}
