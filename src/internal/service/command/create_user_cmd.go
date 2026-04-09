package command

import "github.com/google/uuid"

type CreateUserCmd struct {
	Login      string
	Password   string
	Sex        bool
	Weight     float32
	Age        int
	IsPregnant bool
	IsDriver   bool
	Illnesses  []uuid.UUID
	Allergies  []uuid.UUID
}
