package model

import "time"

type User struct {
	Id        int64     `db:"id"`
	Address   string    `db:"address"`
	CreatedAt time.Time `db:"created_at"`
	OnlineAt  time.Time `db:"online_at"`
}
