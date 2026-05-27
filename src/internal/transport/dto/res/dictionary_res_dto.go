package res

import (
	"github.com/google/uuid"
	"github.com/xtsank/mypills-super-service/src/internal/domain/dictionary"
)

type DictionaryItemDto struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

func NewDictionaryItems(items []dictionary.Item) []DictionaryItemDto {
	res := make([]DictionaryItemDto, 0, len(items))
	for _, item := range items {
		res = append(res, DictionaryItemDto{ID: item.ID, Name: item.Name})
	}
	return res
}

