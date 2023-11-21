package model

type Chain struct {
	Id      int64  `db:"id"`
	ChainId int64  `db:"chain_id"`
	Name    string `db:"name"`
}
