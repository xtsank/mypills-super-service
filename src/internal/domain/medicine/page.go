package medicine

type Page struct {
	Items  []*Medicine
	Total  int
	Limit  int
	Offset int
}
