package model

import "time"

type User struct {
	ID        int64
	Name      string
	Email     string
	Roles     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserCreate struct {
	Name     string
	Password string
	Email    string
	Roles    string
}

type UserPatch struct {
	ID    int64
	Name  *string
	Email *string
	Roles *string
}
