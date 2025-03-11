package responses

import (
	"draft-notification/models"
)

type GetAllUserDeliveryServerResponse struct {
	List       []models.UserDeliveryServer `json:"list"`
	Pagination Pagination                  `json:"pagination"`
}
