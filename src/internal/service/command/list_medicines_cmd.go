package command

import "github.com/google/uuid"

type ListMedicinesCmd struct {
	ID     *uuid.UUID
	Name   *string
	Limit  int
	Offset int
}

func NewListMedicinesCmd(id *uuid.UUID, name *string, limit, offset int) *ListMedicinesCmd {
	return &ListMedicinesCmd{ID: id, Name: name, Limit: limit, Offset: offset}
}
