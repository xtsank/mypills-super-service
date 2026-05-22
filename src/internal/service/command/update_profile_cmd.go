package command

import "github.com/google/uuid"

type UpdateProfileCmd struct {
	ID         uuid.UUID
	Email      *string
	Sex        *bool
	Weight     *int
	Age        *int
	IsPregnant *bool
	IsDriver   *bool
	NotifyEnabled         *bool
	NotifyIntervalMinutes *int
	Illnesses  []uuid.UUID
	Allergies  []uuid.UUID
}

func NewUpdateProfileCmd(id uuid.UUID,
	email *string,
	sex *bool,
	weight *int,
	age *int,
	isPregnant *bool,
	isDriver *bool,
	notifyEnabled *bool,
	notifyIntervalMinutes *int,
	illnesses []uuid.UUID,
	allergies []uuid.UUID,
) *UpdateProfileCmd {
	return &UpdateProfileCmd{
		ID:         id,
		Email:      email,
		Sex:        sex,
		Weight:     weight,
		Age:        age,
		IsPregnant: isPregnant,
		IsDriver:   isDriver,
		NotifyEnabled:         notifyEnabled,
		NotifyIntervalMinutes: notifyIntervalMinutes,
		Illnesses:  illnesses,
		Allergies:  allergies,
	}
}
