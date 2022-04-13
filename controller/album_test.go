package controller_test

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

	test_cases := []struct {
		name   string
		status int
	}{
		{
			name:   "get all albums",
			status: http.StatusOK,
		},
	}

	for _, tc := range test_cases {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", apiprefix+"/albums", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, tc.status, w.Code)
	}
}

func TestPostAlbumRoute(t *testing.T) {

	test_cases := []struct {
		name     string
		body     []byte
		token    string
		response gin.H
		status   int
	}{
		{
			name: "create an album successfully",
			body: []byte(`{
				"title": "New album",
				"artist": "Me Owais",
				"price": 10
			}`),
			token:  "owais",
			status: http.StatusOK,
		},
		{
			name: "try to create album with wrong token",
			body: []byte(`{
				"title": "New album",
				"artist": "Me Owais",
				"price": 10
			}`),
			token:    "wrong_token",
			response: gin.H{"message": "wrong token"},
			status:   http.StatusUnprocessableEntity,
		},
		{
			name: "try to create album with invalid body",
			body: []byte(`{
				"title": "New album",
				"artist": "Me Owais",
			}`),
			token:    "owais",
			response: gin.H{"message": "invalid data"},
			status:   http.StatusUnprocessableEntity,
		},
	}

	for _, tc := range test_cases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", apiprefix+"/albums", bytes.NewBuffer(tc.body))
			req.Header.Add("Authorization", "Bearer "+tc.token)

			router.ServeHTTP(w, req)

			var postRes PostResponse
			json.Unmarshal(w.Body.Bytes(), &postRes)

			if tc.status == http.StatusOK {
				// storing 'id' in memory to use in other test cases
				id = postRes.InsertedID
			}

			expectedResBody, _ := json.Marshal(tc.response)
			resBody := bytes.NewBuffer(expectedResBody).String()

			assert.Equal(t, tc.status, w.Code)

			if tc.status != http.StatusOK {
				assert.Equal(t, resBody, fmt.Sprint(w.Body.String()))
			}
		})
	}
}

func TestGetAlbumByIDRoute(t *testing.T) {

	test_cases := []struct {
		name     string
		id       string
		response gin.H
		status   int
	}{
		{
			name:   "get an album by id",
			id:     id,
			status: http.StatusOK,
		},
		{
			name:     "get an album by wrong id",
			id:       "1",
			response: gin.H{"error": "album not found"},
			status:   http.StatusNotFound,
		},
	}

	for _, tc := range test_cases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", apiprefix+"/albums/"+tc.id, nil)
			router.ServeHTTP(w, req)

			var album models.Album
			json.Unmarshal(w.Body.Bytes(), &album)

			assert.Equal(t, tc.status, w.Code)

			if tc.status == http.StatusOK {
				assert.Equal(t, tc.id, album.ID.Hex())
			} else {
				expectedResBody, _ := json.Marshal(tc.response)
				resBody := bytes.NewBuffer(expectedResBody).String()

				assert.Equal(t, resBody, fmt.Sprint(w.Body.String()))
			}
		})
	}
}

func TestPatchAlbumRoute(t *testing.T) {

	test_cases := []struct {
		name     string
		id       string
		body     []byte
		response gin.H
		status   int
	}{
		{
			name: "update an album",
			id:   id,
			body: []byte(`{
				"title": "New album",
				"price": 10
			}`),
			response: gin.H{"message": "successfully updated the album"},
			status:   http.StatusOK,
		},
		{
			name: "update an album by wrong id",
			id:   "1",
			body: []byte(`{
				"title": "New album",
				"price": 10
			}`),
			response: gin.H{"error": "album not found"},
			status:   http.StatusNotFound,
		},
	}

	for _, tc := range test_cases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("PATCH", apiprefix+"/albums/"+tc.id, bytes.NewBuffer(tc.body))
			router.ServeHTTP(w, req)

			expectedResBody, _ := json.Marshal(tc.response)
			resBody := bytes.NewBuffer(expectedResBody).String()

			assert.Equal(t, tc.status, w.Code)
			assert.Equal(t, resBody, fmt.Sprint(w.Body.String()))
		})
	}
}

func TestDeleteAlbumRoute(t *testing.T) {

	test_cases := []struct {
		name     string
		id       string
		body     []byte
		response gin.H
		status   int
	}{
		{
			name: "update an album",
			id:   id,
			body: []byte(`{
				"title": "New album",
				"price": 10
			}`),
			response: gin.H{"message": "successfully deleted the album"},
			status:   http.StatusOK,
		},
		{
			name: "update an album by wrong id",
			id:   "1",
			body: []byte(`{
				"title": "New album",
				"price": 10
			}`),
			response: gin.H{"error": "album not found"},
			status:   http.StatusNotFound,
		},
	}

	for _, tc := range test_cases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("DELETE", apiprefix+"/albums/"+tc.id, nil)
			router.ServeHTTP(w, req)

			expectedResBody, _ := json.Marshal(tc.response)
			resBody := bytes.NewBuffer(expectedResBody).String()

			assert.Equal(t, tc.status, w.Code)
			assert.Equal(t, resBody, fmt.Sprint(w.Body.String()))
		})
	}
}
