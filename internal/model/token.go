package model

import (
	"github.com/ethereum/go-ethereum/common"
	"time"
)

type Token struct {
	Id        int64          `db:"id"`
	ChainId   int64          `db:"chain_id"`
	Address   common.Address `db:"address"`
	Name      string         `db:"name"`
	Symbol    string         `db:"symbol"`
	Decimals  int            `db:"decimals"`
	CreatedAt time.Time      `db:"created_at"`
}
