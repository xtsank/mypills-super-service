package command

import "github.com/google/uuid"

type UpdateProfileCmd struct {
	ID         uuid.UUID
	Sex        *bool
	Weight     *int
	Age        *int
	IsPregnant *bool
	IsDriver   *bool
	Illnesses  []uuid.UUID
	Allergies  []uuid.UUID
}
