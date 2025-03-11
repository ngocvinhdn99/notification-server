package responses

import (
	"draft-notification/models"
)

type GetAllWebviewServerResponse struct {
	List       []models.WebviewServer `json:"list"`
	Pagination Pagination             `json:"pagination"`
}
