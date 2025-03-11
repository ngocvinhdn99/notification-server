package responses

import (
	"draft-notification/models"
)

type GetAllConnectionResponse struct {
	List       []models.ConnectionResponse `json:"list"`
	Pagination Pagination                  `json:"pagination"`
}
