package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// album represents data about a record album.
type Album struct {
	ID         primitive.ObjectID `bson:"_id" json:"_id"`
	Title      string             `json:"title" validate:"required"`
	Artist     string             `json:"artist" validate:"required"`
	Price      float64            `json:"price" validate:"required"`
	Created_at time.Time          `json:"created_at"`
	Updated_at time.Time          `json:"updated_at"`
}

type AddAlbum struct {
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}
