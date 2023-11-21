package service

import (
	"api-service/internal/model"
	"context"
	"math/big"
)

type TokenRepository interface {
	GetTokensByChainId(ctx context.Context, chainId *big.Int) ([]model.Token, error)
	GetCoinByChainId(ctx context.Context, chainId *big.Int) (model.Token, error)
}

type Token struct {
	tokenRepo TokenRepository
}

func NewToken(repo TokenRepository) *Token {
	return &Token{
		tokenRepo: repo,
	}
}

func (s *Token) Tokens(ctx context.Context, chainId *big.Int) ([]model.Token, error) {
	tokens, err := s.tokenRepo.GetTokensByChainId(ctx, chainId)
	if err != nil {
		return nil, err
	}

	return tokens, nil
}

func (s *Token) Coin(ctx context.Context, chainId *big.Int) (model.Token, error) {
	coin, err := s.tokenRepo.GetCoinByChainId(ctx, chainId)
	if err != nil {
		return model.Token{}, err
	}

	return coin, nil
}
