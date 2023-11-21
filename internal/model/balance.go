package model

import (
	"math/big"
)

type Balance struct {
	Chain        string
	Balance      *big.Int
	Name         string
	Symbol       string
	Decimals     int
	ContractAddr string
}
