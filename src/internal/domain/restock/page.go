package restock

type Page struct {
	Items  []Item
	Total  int
	Limit  int
	Offset int
}
