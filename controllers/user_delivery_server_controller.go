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

var userDeliveryServerCollection *mongo.Collection = configs.GetCollection(configs.DB, "user-delivery-server")

// Helper function for binding and validating request body
func bindAndValidateUserDeliveryServer(c echo.Context, userDeliveryServer *models.UserDeliveryServer) error {
	// Bind request body to userDeliveryServer struct
	if err := c.Bind(userDeliveryServer); err != nil {
		return err
	}

	// Validate the userDeliveryServer struct
	if validationErr := helpers.Validate.Struct(userDeliveryServer); validationErr != nil {
		return validationErr
	}

	return nil
}

// Helper function to handle error response

func CreateUserDeliveryServer(c echo.Context) error {
	ctx, cancel := helpers.CreateContext()
	defer cancel()

	var userDeliveryServer models.UserDeliveryServer
	var findUserDeliveryServer models.UserDeliveryServer

	// Bind and validate the request body
	if err := bindAndValidateUserDeliveryServer(c, &userDeliveryServer); err != nil {
		return helpers.HandleError(c, http.StatusBadRequest, err.Error())
	}

	existNameErr := userDeliveryServerCollection.FindOne(ctx, bson.M{"name": userDeliveryServer.Name}).Decode(&findUserDeliveryServer)
	if existNameErr == nil {
		return helpers.HandleError(c, http.StatusInternalServerError, "Tên user delivery server đã tồn tại")
	}

	// Create new user delivery server
	newUserDeliveryServer := models.UserDeliveryServer{
		Id:        primitive.NewObjectID(),
		Name:      userDeliveryServer.Name,
		Status:    "inactive",
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	result, err := userDeliveryServerCollection.InsertOne(ctx, newUserDeliveryServer)
	if err != nil {
		return helpers.HandleError(c, http.StatusInternalServerError, err.Error())
	}

	return helpers.HandleSuccess(c, result)
}

func GetAllUserDeliveryServers(c echo.Context) error {
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

	var userDeliveryServers []models.UserDeliveryServer
	filter := bson.M{}

	if keyword != "" {
		filter["name"] = bson.M{"$regex": keyword, "$options": "i"}
	}
	if status == "active" || status == "inactive" {
		filter["status"] = status
	}

	results, err := userDeliveryServerCollection.Find(ctx, filter, options.Find().SetLimit(int64(limit)).SetSkip(int64(page*limit)))
	if err != nil {
		return helpers.HandleError(c, http.StatusInternalServerError, err.Error())
	}
	defer results.Close(ctx)

	if err := results.All(ctx, &userDeliveryServers); err != nil {
		return helpers.HandleError(c, http.StatusInternalServerError, err.Error())
	}

	if len(userDeliveryServers) == 0 {
		userDeliveryServers = []models.UserDeliveryServer{}
	}

	// Count total documents matching the filter
	totalCount, err := userDeliveryServerCollection.CountDocuments(ctx, filter)
	if err != nil {
		return helpers.HandleError(c, http.StatusInternalServerError, err.Error())
	}

	// Prepare the response data
	data := responses.GetAllUserDeliveryServerResponse{
		List: userDeliveryServers,
		Pagination: responses.Pagination{
			Total: int(totalCount),
			Limit: limit,
			Page:  page,
		}}
	return helpers.HandleSuccess(c, data)
}

func GetUserDeliveryServerDetail(c echo.Context) error {
	ctx, cancel := helpers.CreateContext()
	defer cancel()

	id := c.Param("id")
	objId, _ := primitive.ObjectIDFromHex(id)

	var userDeliveryServer models.UserDeliveryServer
	err := userDeliveryServerCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&userDeliveryServer)
	if err != nil {
		return helpers.HandleError(c, http.StatusInternalServerError, err.Error())
	}

	return helpers.HandleSuccess(c, userDeliveryServer)
}

func UpdateUserDeliveryServer(c echo.Context) error {
	ctx, cancel := helpers.CreateContext()
	defer cancel()

	id := c.Param("id")
	objId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return helpers.HandleError(c, http.StatusBadRequest, "Id không đúng định dạng")
	}

	var userDeliveryServer models.UserDeliveryServer
	var findUserDeliveryServer models.UserDeliveryServer

	// Bind and validate the request body
	if err := bindAndValidateUserDeliveryServer(c, &userDeliveryServer); err != nil {
		return helpers.HandleError(c, http.StatusBadRequest, err.Error())
	}

	fmt.Println(userDeliveryServer)

	existIdErr := userDeliveryServerCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&findUserDeliveryServer)
	if existIdErr != nil {
		return helpers.HandleError(c, http.StatusInternalServerError, "ID không tồn tại trong DB")
	}
	if findUserDeliveryServer.Name == userDeliveryServer.Name {
		return helpers.HandleSuccess(c, "thành công")
	}

	fmt.Println(userDeliveryServer.Name)

	existNameErr := userDeliveryServerCollection.FindOne(ctx, bson.M{"name": userDeliveryServer.Name}).Decode(&findUserDeliveryServer)
	if existNameErr == nil {
		return helpers.HandleError(c, http.StatusInternalServerError, "Tên user delivery server đã tồn tại")
	}

	update := bson.M{"name": userDeliveryServer.Name}
	result, err := userDeliveryServerCollection.UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$set": update})
	if err != nil {
		return helpers.HandleError(c, http.StatusInternalServerError, err.Error())
	}

	// Retrieve updated user delivery server details
	var updatedUserDeliveryServer models.UserDeliveryServer
	if result.MatchedCount == 1 {
		err := userDeliveryServerCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&updatedUserDeliveryServer)
		if err != nil {
			return helpers.HandleError(c, http.StatusInternalServerError, err.Error())
		}
	}

	return helpers.HandleSuccess(c, updatedUserDeliveryServer)
}

func ChangeStatusUserDeliveryServer(c echo.Context) error {
	ctx, cancel := helpers.CreateContext()
	defer cancel()

	id := c.Param("id")
	objId, _ := primitive.ObjectIDFromHex(id)

	var userDeliveryServer models.UserDeliveryServer
	err := userDeliveryServerCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&userDeliveryServer)
	if err != nil {
		return helpers.HandleError(c, http.StatusInternalServerError, "Id ko tồn tại trong DB")
	}

	var request dtos.ChangeStatusUserDeliveryServerRequest
	if err := c.Bind(&request); err != nil {
		return helpers.HandleError(c, http.StatusBadRequest, "Invalid JSON format")
	}

	if request.Status != "active" && request.Status != "inactive" {
		return helpers.HandleError(c, http.StatusBadRequest, "Invalid status value")
	}

	if request.Status == userDeliveryServer.Status {
		return helpers.HandleSuccess(c, "thành công")
	}

	update := bson.M{"status": request.Status}
	result, err := userDeliveryServerCollection.UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$set": update})
	if err != nil {
		return helpers.HandleError(c, http.StatusInternalServerError, err.Error())
	}

	if request.Status == "inactive" {
		results, err := connectionCollection.Find(ctx, bson.M{"userdeliveryserverid": userDeliveryServer.Id})

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

	// Retrieve updated user delivery server details
	var updatedUserDeliveryServer models.UserDeliveryServer
	if result.MatchedCount == 1 {
		err := userDeliveryServerCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&updatedUserDeliveryServer)
		if err != nil {
			return helpers.HandleError(c, http.StatusInternalServerError, err.Error())
		}
	}

	return helpers.HandleSuccess(c, updatedUserDeliveryServer)
}
