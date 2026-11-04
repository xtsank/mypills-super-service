package req

import "github.com/google/uuid"

type ProfilePatchDto struct {
	Email      *string     `json:"email"`
	Age        *int        `json:"age"`
	Weight     *int        `json:"weight"`
	Sex        *bool       `json:"sex"`
	IsPregnant *bool       `json:"is_pregnant"`
	IsDriver   *bool       `json:"is_driver"`
	Allergies  []uuid.UUID `json:"allergies"`
	Illnesses  []uuid.UUID `json:"illnesses"`
}
