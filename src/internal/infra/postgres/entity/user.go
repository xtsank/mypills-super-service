package entity

import (
	"time"

	"github.com/google/uuid"
)

type UserEntity struct {
	ID         uuid.UUID  `db:"id"`
	Login      string     `db:"login"`
	Email      string     `db:"email"`
	Password   string     `db:"password"`
	IsAdmin    bool       `db:"is_admin"`
	Sex        bool       `db:"sex"`
	Weight     int        `db:"weight"`
	Age        int        `db:"age"`
	IsPregnant bool       `db:"is_pregnant"`
	IsDriver   bool       `db:"is_driver"`
	NotifyEnabled        bool       `db:"notify_enabled"`
	NotifyIntervalMinutes int       `db:"notify_interval_minutes"`
	LastNotifiedAt       *time.Time `db:"last_notified_at"`
}
