package user

import "github.com/google/uuid"

type User struct {
	ID         uuid.UUID
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
