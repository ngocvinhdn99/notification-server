package controllers

import (
	"draft-notification/configs"
	"draft-notification/dtos"
	"draft-notification/helpers"
	"draft-notification/models"
	"draft-notification/responses"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var connectionCollection *mongo.Collection = configs.GetCollection(configs.DB, "connection")

// Helper function for binding and validating request body
func bindAndValidateConnection(c echo.Context, connection *models.Connection) error {
	// Bind request body to connection struct
	if err := c.Bind(connection); err != nil {
		return err
	}

	// Validate the connection struct
	if validationErr := helpers.Validate.Struct(connection); validationErr != nil {
		return validationErr
	}

	return nil
}

func CreateConnection(c echo.Context) error {
	ctx, cancel := helpers.CreateContext()
	defer cancel()

	userDeliveryServerId := c.Param("userDeliveryServerId")

	var connection models.Connection

	// Bind and validate the request body
	if err := bindAndValidateConnection(c, &connection); err != nil {
		return helpers.HandleError(c, http.StatusBadRequest, err.Error())
	}

	userDeliveryServerObjId, err := primitive.ObjectIDFromHex(userDeliveryServerId)
	if err != nil {
		return helpers.HandleError(c, http.StatusBadRequest, "Invalid UserDeliveryServerId")
	}

	if count, _ := userDeliveryServerCollection.CountDocuments(ctx, bson.M{"_id": userDeliveryServerObjId}); count == 0 {
		return helpers.HandleError(c, http.StatusInternalServerError, "ID User delivery server không tồn tại trong DB")
	}

	if count, _ := webviewServerCollection.CountDocuments(ctx, bson.M{"_id": connection.WebviewServerId}); count == 0 {
		return helpers.HandleError(c, http.StatusInternalServerError, "ID webview server không tồn tại trong DB")
	}

	if connection.UserDeliveryServerWebHookUrl != "" {
		if _, err := url.ParseRequestURI(connection.UserDeliveryServerWebHookUrl); err != nil {
			return helpers.HandleError(c, http.StatusBadRequest, "Invalid UserDeliveryServerWebHookUrl")
		}
	}

	count, _ := connectionCollection.CountDocuments(ctx, bson.M{"userdeliveryserverid": userDeliveryServerObjId,
		"webviewserverid": connection.WebviewServerId})

	if count > 0 {
		return helpers.HandleError(c, http.StatusInternalServerError, "Connection giữa webview server và user delivery server đã tồn tại")
	}

	WebviewServerApiKey, err := helpers.GenerateAPIKey(32)
	if err != nil {
		panic(err)
	}

	UserDeliveryServerApiKey, err := helpers.GenerateAPIKey(32)
	if err != nil {
		panic(err)
	}

	// Create new connection
	newConnection := models.Connection{
		Id:                           primitive.NewObjectID(),
		CreatedAt:                    time.Now().UTC(),
		UpdatedAt:                    time.Now().UTC(),
		Status:                       "inactive",
		WebviewServerApiKey:          WebviewServerApiKey,
		UserDeliveryServerApiKey:     UserDeliveryServerApiKey,
		WebviewServerId:              connection.WebviewServerId,
		UserDeliveryServerId:         userDeliveryServerObjId,
		UserDeliveryServerWebHookUrl: connection.UserDeliveryServerWebHookUrl,
	}

	result, err := connectionCollection.InsertOne(ctx, newConnection)
	if err != nil {
		return helpers.HandleError(c, http.StatusInternalServerError, err.Error())
	}

	return helpers.HandleSuccess(c, result)
}

func GetAllConnections(c echo.Context) error {
	ctx, cancel := helpers.CreateContext()
	defer cancel()

	userDeliveryServerId := c.Param("userDeliveryServerId")

	userDeliveryServerObjId, err := primitive.ObjectIDFromHex(userDeliveryServerId)
	if err != nil {
		return helpers.HandleError(c, http.StatusBadRequest, "Invalid UserDeliveryServerId")
	}

	if count, _ := userDeliveryServerCollection.CountDocuments(ctx, bson.M{"_id": userDeliveryServerObjId}); count == 0 {
		return helpers.HandleError(c, http.StatusInternalServerError, "ID User delivery server không tồn tại trong DB")
	}

	webviewServerId := c.QueryParam("webviewServerId")
	status := c.QueryParam("status")
	limitStr := c.QueryParam("limit")
	pageStr := c.QueryParam("page")

	limit := 10
	page := 0

	if limitStr != "" {
		limitParsed, err := strconv.Atoi(limitStr)
		if err == nil && limitParsed > 0 {
			limit = limitParsed
		}
	}

	if pageStr != "" {
		parsedPage, err := strconv.Atoi(pageStr)
		if err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}

	filter := bson.M{
		"userdeliveryserverid": userDeliveryServerObjId,
	}

	if webviewServerObjId, err := primitive.ObjectIDFromHex(webviewServerId); err == nil {
		filter["webviewserverid"] = webviewServerObjId
	}

	if status == "active" || status == "inactive" {
		filter["status"] = status
	}

	results, err := connectionCollection.Find(ctx, filter, options.Find().SetLimit(int64(limit)).SetSkip(int64(page*limit)))
	if err != nil {
		return helpers.HandleError(c, http.StatusInternalServerError, err.Error())
	}
	defer results.Close(ctx)

	// Tạo danh sách chứa dữ liệu phản hồi
	var connectionResponses []models.ConnectionResponse

	for results.Next(ctx) {
		var conn models.Connection

		if err := results.Decode(&conn); err != nil {
			return helpers.HandleError(c, http.StatusInternalServerError, err.Error())
		}

		// Khai báo struct chứa server info
		var webviewServer models.ServerInfo
		var userDeliveryServer models.ServerInfo

		// Lấy thông tin từ collection webviewServer
		if err := webviewServerCollection.FindOne(ctx, bson.M{"_id": conn.WebviewServerId}).Decode(&webviewServer); err != nil {
			// Nếu không tìm thấy, giữ nguyên giá trị rỗng
			webviewServer = models.ServerInfo{}
		}

		// Lấy thông tin từ collection userDeliveryServer
		if err := userDeliveryServerCollection.FindOne(ctx, bson.M{"_id": conn.UserDeliveryServerId}).Decode(&userDeliveryServer); err != nil {
			userDeliveryServer = models.ServerInfo{}
		}

		// Gộp dữ liệu vào struct response
		connectionResponse := models.ConnectionResponse{
			Id:                           conn.Id,
			Status:                       conn.Status,
			CreatedAt:                    conn.CreatedAt,
			UpdatedAt:                    conn.UpdatedAt,
			WebviewServerApiKey:          conn.WebviewServerApiKey,
			UserDeliveryServerApiKey:     conn.UserDeliveryServerApiKey,
			WebviewServer:                webviewServer,
			UserDeliveryServer:           userDeliveryServer,
			UserDeliveryServerWebHookUrl: conn.UserDeliveryServerWebHookUrl,
		}

		connectionResponses = append(connectionResponses, connectionResponse)
	}

	// Kiểm tra nếu không có dữ liệu thì trả về mảng rỗng
	if len(connectionResponses) == 0 {
		connectionResponses = []models.ConnectionResponse{}
	}

	// Count total documents matching the filter
	totalCount, err := connectionCollection.CountDocuments(ctx, filter)
	if err != nil {
		return helpers.HandleError(c, http.StatusInternalServerError, err.Error())
	}

	// Prepare the response data
	data := responses.GetAllConnectionResponse{
		List: connectionResponses,
		Pagination: responses.Pagination{
			Total: int(totalCount),
			Limit: limit,
			Page:  page,
		}}
	return helpers.HandleSuccess(c, data)
}

func UpdateConnectionWebhookUrl(c echo.Context) error {
	ctx, cancel := helpers.CreateContext()
	defer cancel()

	id := c.Param("id")
	objId, objIdErr := primitive.ObjectIDFromHex(id)

	if objIdErr != nil {
		return helpers.HandleError(c, http.StatusBadRequest, "Invalid ID")
	}

	var connection models.Connection
	err := connectionCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&connection)
	if err != nil {
		return helpers.HandleError(c, http.StatusInternalServerError, "Id ko tồn tại trong DB")
	}

	var request dtos.UpdateConnectionWebhookUrlRequest
	if err := c.Bind(&request); err != nil {
		return helpers.HandleError(c, http.StatusBadRequest, "Invalid JSON format")
	}

	if _, err := url.ParseRequestURI(request.UserDeliveryServerWebHookUrl); err != nil {
		return helpers.HandleError(c, http.StatusBadRequest, "Invalid UserDeliveryServerWebHookUrl")
	}

	if request.UserDeliveryServerWebHookUrl == connection.UserDeliveryServerWebHookUrl {
		return helpers.HandleSuccess(c, "thành công")
	}

	update := bson.M{"userdeliveryserverwebhookurl": request.UserDeliveryServerWebHookUrl}
	result, err := connectionCollection.UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$set": update})
	if err != nil {
		return helpers.HandleError(c, http.StatusInternalServerError, err.Error())
	}

	// Retrieve updated connection details
	var updatedConnection models.Connection
	if result.MatchedCount == 1 {
		err := connectionCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&updatedConnection)
		if err != nil {
			return helpers.HandleError(c, http.StatusInternalServerError, err.Error())
		}
	}

	return helpers.HandleSuccess(c, updatedConnection)
}

func ChangeStatusConnection(c echo.Context) error {
	ctx, cancel := helpers.CreateContext()
	defer cancel()

	id := c.Param("id")
	objId, objIdErr := primitive.ObjectIDFromHex(id)

	if objIdErr != nil {
		return helpers.HandleError(c, http.StatusBadRequest, "Invalid ID")
	}

	var connection models.Connection
	err := connectionCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&connection)
	if err != nil {
		return helpers.HandleError(c, http.StatusInternalServerError, "Id ko tồn tại trong DB")
	}

	var request dtos.ChangeStatusConnectionRequest
	if err := c.Bind(&request); err != nil {
		return helpers.HandleError(c, http.StatusBadRequest, "Invalid JSON format")
	}

	if request.Status != "active" && request.Status != "inactive" {
		return helpers.HandleError(c, http.StatusBadRequest, "Invalid status value")
	}

	if request.Status == "active" {
		var webviewServer models.WebviewServer
		var userDeliveryServer models.UserDeliveryServer

		if err := webviewServerCollection.FindOne(ctx, bson.M{"_id": connection.WebviewServerId}).Decode(&webviewServer); err != nil {
			return helpers.HandleError(c, http.StatusInternalServerError, "Không tìm thấy thông tin webview server")
		}

		if err := userDeliveryServerCollection.FindOne(ctx, bson.M{"_id": connection.UserDeliveryServerId}).Decode(&userDeliveryServer); err != nil {
			return helpers.HandleError(c, http.StatusInternalServerError, "Không tìm thấy thông tin user delivery server")
		}

		if userDeliveryServer.Status != "active" {
			return helpers.HandleError(c, http.StatusInternalServerError, "User delivery server chưa active")
		}

		if webviewServer.Status != "active" {
			return helpers.HandleError(c, http.StatusInternalServerError, "Webview server chưa active")
		}
	}

	if request.Status == connection.Status {
		return helpers.HandleSuccess(c, "thành công")
	}

	update := bson.M{"status": request.Status}
	result, err := connectionCollection.UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$set": update})
	if err != nil {
		return helpers.HandleError(c, http.StatusInternalServerError, err.Error())
	}

	// Retrieve updated connection details
	var updatedConnection models.Connection
	if result.MatchedCount == 1 {
		err := connectionCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&updatedConnection)
		if err != nil {
			return helpers.HandleError(c, http.StatusInternalServerError, err.Error())
		}
	}

	return helpers.HandleSuccess(c, updatedConnection)
}
