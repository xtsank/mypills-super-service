package req

import "github.com/google/uuid"

type UpdateProfileDto struct {
	Email                 *string    `json:"email"`
	Sex                   *bool      `json:"sex"`
	Weight                *int       `json:"weight"`
	Age                   *int       `json:"age"`
	IsPregnant            *bool      `json:"is_pregnant"`
	IsDriver              *bool      `json:"is_driver"`
	NotifyEnabled         *bool      `json:"notify_enabled"`
	NotifyIntervalMinutes *int       `json:"notify_interval_minutes"`
	Illnesses             []uuid.UUID `json:"illnesses"`
	Allergies             []uuid.UUID `json:"allergies"`
}
