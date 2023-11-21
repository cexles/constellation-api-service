package response

import "math/big"

type UserBalance struct {
	Name         string               `json:"name"`
	Symbol       string               `json:"symbol"`
	Decimals     int                  `json:"decimals"`
	TokenBalance *big.Int             `json:"tokenBalance"`
	Price        float64              `json:"price"`
	Chains       map[string]ChainInfo `json:"chains"`
}

type ChainInfo struct {
	TokenBalance *big.Int `json:"tokenBalance"`
}
