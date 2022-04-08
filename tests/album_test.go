package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"rest/models"
	"rest/routes"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var router = routes.Routes()

var id string

var apiprefix = "/api/v1"

type PostResponse struct {
	InsertedID string
}

func TestGetAlumbsRoute(t *testing.T) {

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", apiprefix+"/albums", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPostAlbumRoute(t *testing.T) {

	body := []byte(`{
		"title": "New album",
		"artist": "Me Owais",
		"price": 10
	}`)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", apiprefix+"/albums", bytes.NewBuffer(body))
	req.Header.Add("Authorization", "Bearer owais")

	router.ServeHTTP(w, req)

	var postRes PostResponse
	json.Unmarshal(w.Body.Bytes(), &postRes)

	id = postRes.InsertedID
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetAlbumByIDRoute(t *testing.T) {

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", apiprefix+"/albums/"+id, nil)
	router.ServeHTTP(w, req)

	var album models.Album
	json.Unmarshal(w.Body.Bytes(), &album)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, id, album.ID.Hex())
}

func TestGetAlbumByIDRoute_NotFound(t *testing.T) {

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", apiprefix+"/albums/1", nil)
	router.ServeHTTP(w, req)

	expectedResBody, _ := json.Marshal(gin.H{"error": "album not found"})
	resBody := bytes.NewBuffer(expectedResBody).String()

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, resBody, fmt.Sprint(w.Body.String()))
}

func TestPostAlbumRoute_InvalidToken(t *testing.T) {

	body := []byte(`{
		"title": "New album",
		"artist": "Me, Owais",
		"price": 100
	}`)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", apiprefix+"/albums", bytes.NewBuffer(body))
	req.Header.Add("Authorization", "Bearer wrong_token")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

func TestPatchAlbumRoute(t *testing.T) {

	body := []byte(`{
		"title": "New album",
		"price": 10
	}`)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", apiprefix+"/albums/"+id, bytes.NewBuffer(body))

	router.ServeHTTP(w, req)

	expectedResBody, _ := json.Marshal(gin.H{"message": "successfully updated the album"})
	resBody := bytes.NewBuffer(expectedResBody).String()

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, resBody, fmt.Sprint(w.Body.String()))
}

func TestPatchAlbumRoute_NotFound(t *testing.T) {

	body := []byte(`{
		"title": "New album",
		"price": 10
	}`)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", apiprefix+"/albums/1", bytes.NewBuffer(body))

	router.ServeHTTP(w, req)

	expectedResBody, _ := json.Marshal(gin.H{"error": "album not found"})
	resBody := bytes.NewBuffer(expectedResBody).String()

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, resBody, fmt.Sprint(w.Body.String()))
}

func TestDeleteAlbumRoute(t *testing.T) {

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", apiprefix+"/albums/"+id, nil)

	router.ServeHTTP(w, req)

	expectedResBody, _ := json.Marshal(gin.H{"message": "successfully deleted the album"})
	resBody := bytes.NewBuffer(expectedResBody).String()

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, resBody, fmt.Sprint(w.Body.String()))
}

func TestDeleteAlbumRoute_NotFound(t *testing.T) {

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", apiprefix+"/albums/1", nil)

	router.ServeHTTP(w, req)

	expectedResBody, _ := json.Marshal(gin.H{"error": "album not found"})
	resBody := bytes.NewBuffer(expectedResBody).String()

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, resBody, fmt.Sprint(w.Body.String()))
}
