package res

type ReferenceIndexDto struct {
	Illnesses  string `json:"illnesses"`
	Substances string `json:"substances"`
	Forms      string `json:"forms"`
	Units      string `json:"units"`
}

type ReferencePageDto struct {
	Items  []DictionaryItemDto `json:"items"`
	Total  int                 `json:"total"`
	Limit  int                 `json:"limit"`
	Offset int                 `json:"offset"`
}

func NewReferenceIndexDto() *ReferenceIndexDto {
	return &ReferenceIndexDto{
		Illnesses: "/api/v2/reference/illnesses", Substances: "/api/v2/reference/substances",
		Forms: "/api/v2/reference/forms", Units: "/api/v2/reference/units",
	}
}

func NewReferencePageDto(items []DictionaryItemDto, limit, offset int) *ReferencePageDto {
	page := &ReferencePageDto{Items: []DictionaryItemDto{}, Total: len(items), Limit: limit, Offset: offset}
	start := min(offset, len(items))
	end := min(start+limit, len(items))
	page.Items = append(page.Items, items[start:end]...)
	return page
}
