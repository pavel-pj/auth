package model

import "time"

/*
type UserCreate struct {
	Name     string `db:"name"`
	Email    string `db:"email"`
	Password string `db:"password"`
	Roles    string `db:"role"`
}*/

type User struct {
	ID        int64     `db:"id"`
	Name      string    `db:"name"`
	Email     string    `db:"email"`
	Roles     string    `db:"role"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
