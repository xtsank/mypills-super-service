package cabinet

type ItemPage struct {
	Items  []Item
	Total  int
	Limit  int
	Offset int
}
