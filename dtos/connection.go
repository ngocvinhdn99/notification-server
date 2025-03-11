package dtos

type ChangeStatusConnectionRequest struct {
	Status string `json:"status"`
}

type UpdateConnectionWebhookUrlRequest struct {
	UserDeliveryServerWebHookUrl string `json:"userDeliveryServerWebHookUrl"`
}
