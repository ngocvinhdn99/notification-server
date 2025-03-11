package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Connection struct {
	Id                           primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Status                       string             `json:"status,omitempty"`
	CreatedAt                    time.Time          `json:"createdAt,omitempty"`
	UpdatedAt                    time.Time          `json:"updatedAt,omitempty"`
	WebviewServerApiKey          string             `json:"webviewServerApiKey,omitempty"`
	UserDeliveryServerApiKey     string             `json:"userDeliveryServerApiKey,omitempty"`
	WebviewServerId              primitive.ObjectID `json:"webviewServerId,omitempty" validate:"required"`
	UserDeliveryServerId         primitive.ObjectID `json:"userDeliveryServerId,omitempty"`
	UserDeliveryServerWebHookUrl string             `json:"userDeliveryServerWebHookUrl,omitempty"`
}

// Struct chứa thông tin của WebviewServer & UserDeliveryServer
type ServerInfo struct {
	Id   primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name string             `json:"name,omitempty"`
}

// Struct chứa danh sách connection kèm thông tin server
type ConnectionResponse struct {
	Id                           primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Status                       string             `json:"status,omitempty"`
	CreatedAt                    time.Time          `json:"createdAt,omitempty"`
	UpdatedAt                    time.Time          `json:"updatedAt,omitempty"`
	WebviewServerApiKey          string             `json:"webviewServerApiKey,omitempty"`
	UserDeliveryServerApiKey     string             `json:"userDeliveryServerApiKey,omitempty"`
	WebviewServer                ServerInfo         `json:"webviewServer,omitempty"`
	UserDeliveryServer           ServerInfo         `json:"userDeliveryServer,omitempty"`
	UserDeliveryServerWebHookUrl string             `json:"userDeliveryServerWebHookUrl,omitempty"`
}
