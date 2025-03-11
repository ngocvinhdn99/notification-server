package controllers

import (
	"draft-notification/configs"
	"draft-notification/dtos"
	"draft-notification/helpers"
	"draft-notification/models"
	"draft-notification/responses"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var webviewServerCollection *mongo.Collection = configs.GetCollection(configs.DB, "webview-server")

// Helper function for binding and validating request body
func bindAndValidateWebviewServer(c echo.Context, webviewServer *models.WebviewServer) error {
	// Bind request body to webviewServer struct
	if err := c.Bind(webviewServer); err != nil {
		return err
	}

	// Validate the webviewServer struct
	if validationErr := helpers.Validate.Struct(webviewServer); validationErr != nil {
		return validationErr
	}

	return nil
}

func CreateWebviewServer(c echo.Context) error {
	ctx, cancel := helpers.CreateContext()
	defer cancel()

	var webviewServer models.WebviewServer
	var findWebviewServer models.WebviewServer

	// Bind and validate the request body
	if err := bindAndValidateWebviewServer(c, &webviewServer); err != nil {
		return helpers.HandleError(c, http.StatusBadRequest, err.Error())
	}

	existNameErr := webviewServerCollection.FindOne(ctx, bson.M{"name": webviewServer.Name}).Decode(&findWebviewServer)
	if existNameErr == nil {
		return helpers.HandleError(c, http.StatusInternalServerError, "Tên webview server đã tồn tại")
	}

	// Create new webview server
	newWebviewServer := models.WebviewServer{
		Id:        primitive.NewObjectID(),
		Name:      webviewServer.Name,
		Status:    "inactive",
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	result, err := webviewServerCollection.InsertOne(ctx, newWebviewServer)
	if err != nil {
		return helpers.HandleError(c, http.StatusInternalServerError, err.Error())
	}

	return helpers.HandleSuccess(c, result)
}

func GetAllWebviewServers(c echo.Context) error {
	ctx, cancel := helpers.CreateContext()
	defer cancel()

	keyword := c.QueryParam("keyword")
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

	var webviewServers []models.WebviewServer
	filter := bson.M{}

	if keyword != "" {
		filter["name"] = bson.M{"$regex": keyword, "$options": "i"}
	}
	if status == "active" || status == "inactive" {
		filter["status"] = status
	}

	results, err := webviewServerCollection.Find(ctx, filter, options.Find().SetLimit(int64(limit)).SetSkip(int64(page*limit)))
	if err != nil {
		return helpers.HandleError(c, http.StatusInternalServerError, err.Error())
	}
	defer results.Close(ctx)

	if err := results.All(ctx, &webviewServers); err != nil {
		return helpers.HandleError(c, http.StatusInternalServerError, err.Error())
	}

	if len(webviewServers) == 0 {
		webviewServers = []models.WebviewServer{}
	}

	// Count total documents matching the filter
	totalCount, err := webviewServerCollection.CountDocuments(ctx, filter)
	if err != nil {
		return helpers.HandleError(c, http.StatusInternalServerError, err.Error())
	}

	// Prepare the response data
	data := responses.GetAllWebviewServerResponse{
		List: webviewServers,
		Pagination: responses.Pagination{
			Total: int(totalCount),
			Limit: limit,
			Page:  page,
		}}
	return helpers.HandleSuccess(c, data)
}

func GetWebviewServerDetail(c echo.Context) error {
	ctx, cancel := helpers.CreateContext()
	defer cancel()

	id := c.Param("id")
	objId, _ := primitive.ObjectIDFromHex(id)

	var webviewServer models.WebviewServer
	err := webviewServerCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&webviewServer)
	if err != nil {
		return helpers.HandleError(c, http.StatusInternalServerError, err.Error())
	}

	return helpers.HandleSuccess(c, webviewServer)
}

func UpdateWebviewServer(c echo.Context) error {
	ctx, cancel := helpers.CreateContext()
	defer cancel()

	id := c.Param("id")
	objId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return helpers.HandleError(c, http.StatusBadRequest, "Id không đúng định dạng")
	}

	var webviewServer models.WebviewServer
	var findWebviewServer models.WebviewServer

	// Bind and validate the request body
	if err := bindAndValidateWebviewServer(c, &webviewServer); err != nil {
		return helpers.HandleError(c, http.StatusBadRequest, err.Error())
	}

	fmt.Println(webviewServer)

	existIdErr := webviewServerCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&findWebviewServer)
	if existIdErr != nil {
		return helpers.HandleError(c, http.StatusInternalServerError, "ID không tồn tại trong DB")
	}
	if findWebviewServer.Name == webviewServer.Name {
		return helpers.HandleSuccess(c, "thành công")
	}

	fmt.Println(webviewServer.Name)

	existNameErr := webviewServerCollection.FindOne(ctx, bson.M{"name": webviewServer.Name}).Decode(&findWebviewServer)
	if existNameErr == nil {
		return helpers.HandleError(c, http.StatusInternalServerError, "Tên webview server đã tồn tại")
	}

	update := bson.M{"name": webviewServer.Name}
	result, err := webviewServerCollection.UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$set": update})
	if err != nil {
		return helpers.HandleError(c, http.StatusInternalServerError, err.Error())
	}

	// Retrieve updated webview server details
	var updatedWebviewServer models.WebviewServer
	if result.MatchedCount == 1 {
		err := webviewServerCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&updatedWebviewServer)
		if err != nil {
			return helpers.HandleError(c, http.StatusInternalServerError, err.Error())
		}
	}

	return helpers.HandleSuccess(c, updatedWebviewServer)
}

func ChangeStatusWebviewServer(c echo.Context) error {
	ctx, cancel := helpers.CreateContext()
	defer cancel()

	id := c.Param("id")
	objId, _ := primitive.ObjectIDFromHex(id)

	var webviewServer models.WebviewServer
	err := webviewServerCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&webviewServer)
	if err != nil {
		return helpers.HandleError(c, http.StatusInternalServerError, "Id ko tồn tại trong DB")
	}

	var request dtos.ChangeStatusWebviewServerRequest
	if err := c.Bind(&request); err != nil {
		return helpers.HandleError(c, http.StatusBadRequest, "Invalid JSON format")
	}

	if request.Status != "active" && request.Status != "inactive" {
		return helpers.HandleError(c, http.StatusBadRequest, "Invalid status value")
	}

	if request.Status == webviewServer.Status {
		return helpers.HandleSuccess(c, "thành công")
	}

	update := bson.M{"status": request.Status}
	result, err := webviewServerCollection.UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$set": update})
	if err != nil {
		return helpers.HandleError(c, http.StatusInternalServerError, err.Error())
	}

	if request.Status == "inactive" {
		results, err := connectionCollection.Find(ctx, bson.M{"webviewserverid": webviewServer.Id})

		if err != nil {
			return helpers.HandleError(c, http.StatusInternalServerError, "tìm thông tin connection bị lỗi")
		}

		for results.Next(ctx) {
			var connection models.Connection

			if err := results.Decode(&connection); err != nil {
				return helpers.HandleError(c, http.StatusInternalServerError, err.Error())
			}

			if connection.Status == "active" {
				_, err := connectionCollection.UpdateOne(ctx, bson.M{"_id": connection.Id}, bson.M{"$set": bson.M{"status": "inactive"}})
				if err != nil {
					return helpers.HandleError(c, http.StatusInternalServerError, err.Error())
				}
			}
		}

	}

	// Retrieve updated webview server details
	var updatedWebviewServer models.WebviewServer
	if result.MatchedCount == 1 {
		err := webviewServerCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&updatedWebviewServer)
		if err != nil {
			return helpers.HandleError(c, http.StatusInternalServerError, err.Error())
		}
	}

	return helpers.HandleSuccess(c, updatedWebviewServer)
}
