package request

import (
	"github.com/ethereum/go-ethereum/common"
)

type LoginRequest struct {
	Address   common.Address `json:"address"`
	Signature string         `json:"signature"`
}

type RefreshRequest struct {
	Token string `json:"token"`
}
