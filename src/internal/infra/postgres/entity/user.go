package entity

import "github.com/google/uuid"

type UserEntity struct {
	ID         uuid.UUID `db:"id"`
	Login      string    `db:"login"`
	Password   string    `db:"password"`
	IsAdmin    bool      `db:"is_admin"`
	Sex        bool      `db:"sex"`
	Weight     int       `db:"weight"`
	Age        int       `db:"age"`
	IsPregnant bool      `db:"is_pregnant"`
	IsDriver   bool      `db:"is_driver"`
}
