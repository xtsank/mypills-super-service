package entity

import (
	"time"

	"github.com/google/uuid"
)

type UserEntity struct {
	ID                    uuid.UUID  `bson:"_id"`
	Login                 string     `bson:"login"`
	Email                 string     `bson:"email"`
	Password              string     `bson:"password"`
	IsAdmin               bool       `bson:"is_admin"`
	Sex                   bool       `bson:"sex"`
	Weight                int        `bson:"weight"`
	Age                   int        `bson:"age"`
	IsPregnant            bool       `bson:"is_pregnant"`
	IsDriver              bool       `bson:"is_driver"`
	NotifyEnabled         bool       `bson:"notify_enabled"`
	NotifyIntervalMinutes int        `bson:"notify_interval_minutes"`
	LastNotifiedAt        *time.Time `bson:"last_notified_at,omitempty"`

	Illnesses []uuid.UUID `bson:"illnesses"`
	Allergies []uuid.UUID `bson:"allergies"`
}
