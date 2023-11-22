package service

import (
	"api-service/internal/api/response"
	"api-service/internal/config"
	"api-service/internal/model"
	"context"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
	"time"
)

type UserRepository interface {
	CreateOrUpdate(ctx context.Context, address string) (model.User, error)
}

type Auth struct {
	userRepo UserRepository
	jwtCfg   *config.Jwt
}

func NewAuth(repo UserRepository, jwtCfg *config.Jwt) *Auth {
	return &Auth{
		userRepo: repo,
		jwtCfg:   jwtCfg,
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
		"exp":     time.Now().Add(time.Hour * time.Duration(s.jwtCfg.Expiration)).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(s.jwtCfg.SecretKey))
	if err != nil {
		return response.Login{}, err
	}

	return response.Login{
		Token: t,
	}, nil
}

func (s *Auth) RefreshToken(ctx *fiber.Ctx, token string) (response.RefreshToken, error) {
	refreshTokenString := ctx.Get("Authorization")[7:]

	if !strings.EqualFold(refreshTokenString, token) {
		return response.RefreshToken{}, fiber.NewError(http.StatusBadRequest, "Bad token")
	}

	oldToken, err := jwt.Parse(refreshTokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtCfg.SecretKey), nil
	})
	if err != nil {
		return response.RefreshToken{}, err
	}

	oldClaims, ok := oldToken.Claims.(jwt.MapClaims)
	if !ok {
		return response.RefreshToken{}, err
	}

	address, ok := oldClaims["address"].(string)
	if !ok {
		return response.RefreshToken{}, err
	}

	newClaims := jwt.MapClaims{
		"address": address,
		"exp":     time.Now().Add(time.Hour * time.Duration(s.jwtCfg.Expiration)).Unix(),
	}

	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)
	t, err := newToken.SignedString([]byte(s.jwtCfg.SecretKey))
	if err != nil {
		return response.RefreshToken{}, err
	}

	return response.RefreshToken{
		Token: t,
	}, nil
}
