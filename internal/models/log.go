package models

import "time"

const GetUsersAction = "get users"
const DeleteAction = "delete"
const BlockAction = "block"
const ShowLogAction = "show log"

type Log struct {
	Id        int64     `json:"id"`
	UserId    int64     `json:"user_id"`
	Action    string    `json:"action"`
	CreatedAt time.Time `json:"created_at"`
}
