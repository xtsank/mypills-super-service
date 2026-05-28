package entity

import "github.com/google/uuid"

type DictionaryItemEntity struct {
	ID   uuid.UUID `bson:"_id"`
	Name string    `bson:"name"`
}
