package req

import "github.com/google/uuid"

type UpdateProfileDto struct {
	Sex        *bool       `json:"sex"`
	Weight     *int        `json:"weight"`
	Age        *int        `json:"age"`
	IsPregnant *bool       `json:"is_pregnant"`
	IsDriver   *bool       `json:"is_driver"`
	Illnesses  []uuid.UUID `json:"illnesses"`
	Allergies  []uuid.UUID `json:"allergies"`
}
