package command

type ListReferenceCmd struct {
	Kind   string
	Limit  int
	Offset int
}

func NewListReferenceCmd(kind string, limit, offset int) *ListReferenceCmd {
	return &ListReferenceCmd{Kind: kind, Limit: limit, Offset: offset}
}
