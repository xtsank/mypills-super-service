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

func NewUpdateProfileCmd(id uuid.UUID,
	sex *bool,
	weight *int,
	age *int,
	isPregnant *bool,
	isDriver *bool,
	illnesses []uuid.UUID,
	allergies []uuid.UUID,
) *UpdateProfileCmd {
	return &UpdateProfileCmd{
		ID:         id,
		Sex:        sex,
		Weight:     weight,
		Age:        age,
		IsPregnant: isPregnant,
		IsDriver:   isDriver,
		Illnesses:  illnesses,
		Allergies:  allergies,
	}
}
