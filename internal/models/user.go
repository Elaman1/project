package models

import "time"

type User struct {
	Id        int64     `json:"id"`
	Name      string    `json:"name"`
	Password  string    `json:"password"`
	RoleID    int64     `json:"role_id"`
	Role      Role      `json:"role"`
	Blocked   bool      `json:"blocked"`
	CreatedAt time.Time `json:"created_at"`
}
