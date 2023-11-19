package service

import (
	"api-service/internal/api/response"
	"api-service/internal/model"
	"context"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type UserRepository interface {
	CreateOrUpdate(ctx context.Context, address string) (model.User, error)
}

type Auth struct {
	userRepo UserRepository
}

func NewAuth(repo UserRepository) *Auth {
	return &Auth{
		userRepo: repo,
	}
}

func VerifySignature(address common.Address, signature string) (bool, error) {
	signatureBytes, err := hexutil.Decode(signature)
	if err != nil {
		return false, err
	}

	messageHash := accounts.TextHash([]byte("hello"))

	if signatureBytes[crypto.RecoveryIDOffset] == 27 || signatureBytes[crypto.RecoveryIDOffset] == 28 {
		signatureBytes[crypto.RecoveryIDOffset] -= 27
	}

	publicKey, err := crypto.SigToPub(messageHash, signatureBytes)
	if err != nil {
		return false, err
	}

	signerAddress := crypto.PubkeyToAddress(*publicKey)

	return signerAddress == address, err
}

func (s *Auth) Login(ctx context.Context, address common.Address) (response.Login, error) {
	u, err := s.userRepo.CreateOrUpdate(ctx, address.Hex())
	if err != nil {
		return response.Login{}, err
	}

	claims := jwt.MapClaims{
		"address": u.Address,
		"exp":     time.Now().Add(time.Second * 15).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return response.Login{}, err
	}

	return response.Login{
		Token: t,
	}, nil
}
