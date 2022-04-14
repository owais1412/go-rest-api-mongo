package controller

import (
	"context"
	"log"
	"net/http"
	"rest/database"
	"rest/middlewares"
	"rest/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

//connect to to the database and open an album collection
var client *mongo.Client
var albumsCollection *mongo.Collection

var validate *validator.Validate

func init() {
	client = database.DBinstance()
	albumsCollection = database.OpenCollection(client, "albums")
	validate = validator.New()
}

// GetAlbumByID godoc
// @Summary      Get all albums
// @Description  get albums
// @Tags         albums
// @Accept       json
// @Produce      json
// @Success      200  {array}  	models.Album
// @Failure      400  {object}  models.ErrorMessage
// @Failure      404  {object}  models.ErrorMessage
// @Failure      500  {object}  models.ErrorMessage
// @Router       /albums [get]
func GetAlbums(c *gin.Context) {
	cursor, err := albumsCollection.Find(c, bson.M{})

	if err != nil {
		log.Fatal(err)
	}

	var albums []models.Album

	if err = cursor.All(c, &albums); err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, albums)
}

// GetAlbumByID godoc
// @Summary      Get an album
// @Description  get string by ID
// @Tags         albums
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Album ID"
// @Success      200  {object}  models.Album
// @Failure      400  {object}  models.ErrorMessage
// @Failure      404  {object}  models.ErrorMessage
// @Failure      500  {object}  models.ErrorMessage
// @Router       /albums/{id} [get]
func GetAlbumByID(c *gin.Context) {
	id, _ := primitive.ObjectIDFromHex(c.Param("id"))

	var album models.Album

	err := albumsCollection.FindOne(c, bson.M{"_id": id}).Decode(&album)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "album not found"})
		return
	}

	c.JSON(http.StatusOK, album)
}

// PostAlbum godoc
// @Summary      Add an album
// @Description  add album by json
// @Tags         albums
// @Accept       json
// @Produce      json
// @Param        album   body      models.AddAlbum  true  "Add Album"
// @Success      200	{object}  models.Album
// @Failure      400	{object}  models.ErrorMessage
// @Failure      404	{object}  models.ErrorMessage
// @Failure      500	{object}  models.ErrorMessage
// @Security     bearer
// @Router       /albums [post]
func PostAlbum(c *gin.Context) {
	//this is used to determine how long the API call should last
	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)

	var album models.Album

	token := c.GetHeader("Authorization")

	if !middlewares.IsValidToken(token) {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "wrong token"})
		cancel()
		return
	}

	// Call ShouldBindJSON to bind the received JSON to album.
	if err := c.ShouldBindJSON(&album); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "invalid data"})
		cancel()
		return
	}

	if validationErr := validate.Struct(album); validationErr != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": validationErr.Error()})
		cancel()
		return
	}

	// the RFC3339 layout 2006-01-02T15:04:05Z07:00
	album.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	album.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	//generate new ID for the object to be created
	album.ID = primitive.NewObjectID()

	//insert the newly created object into mongodb
	result, insertErr := albumsCollection.InsertOne(ctx, album)
	if insertErr != nil {
		msg := "Album was not created"
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		cancel()
		return
	}
	defer cancel()

	//return the id of the created object
	c.JSON(http.StatusOK, result)
}

// UpdateAlbum godoc
// @Summary      Update an album
// @Description  Update by json album
// @Tags         albums
// @Accept       json
// @Produce      json
// @Param        id       path      string	true  "Account ID"
// @Param        album	body      models.AddAlbum  true  "Update Album"
// @Success      200      {object}  models.SuccessMessage
// @Failure      400      {object}  models.ErrorMessage
// @Failure      404      {object}  models.ErrorMessage
// @Failure      500      {object}  models.ErrorMessage
// @Router       /albums/{id} [patch]
func UpdateAlbum(c *gin.Context) {
	id, _ := primitive.ObjectIDFromHex(c.Param("id"))

	var album models.Album

	err := albumsCollection.FindOne(c, bson.M{"_id": id}).Decode(&album)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "album not found"})
		return
	}

	// Call ShouldBindJSON to bind the received JSON to album.
	if err = c.ShouldBindJSON(&album); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "invalid data"})
		return
	}

	if validationErr := validate.Struct(&album); validationErr != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": validationErr.Error()})
		return
	}

	album.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	res, _ := albumsCollection.UpdateByID(c, id, bson.M{"$set": album})

	if res.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "album not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "successfully updated the album"})
}

// DeleteAlbumByID godoc
// @Summary      Delete an albums
// @Description  Delete by album ID
// @Tags         albums
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Album ID"
// @Success      200      {object}  models.SuccessMessage
// @Failure      400      {object}  models.ErrorMessage
// @Failure      404      {object}  models.ErrorMessage
// @Failure      500      {object}  models.ErrorMessage
// @Router       /albums/{id} [delete]
func DeleteAlbumByID(c *gin.Context) {
	id, _ := primitive.ObjectIDFromHex(c.Param("id"))

	res, _ := albumsCollection.DeleteOne(c, bson.M{"_id": id})

	if res.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "album not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "successfully deleted the album"})
}
