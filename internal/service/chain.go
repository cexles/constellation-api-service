package service

import (
	"api-service/internal/config"
	"context"
	"math/big"
)

type ChainRepository interface {
	Exist(ctx context.Context, name string) (bool, error)
	GetNameByChainId(ctx context.Context, chainId *big.Int) (string, error)
}

type Chain struct {
	chainRepo ChainRepository
	rpcConfig map[string]*config.RPCDetails
}

func NewChain(repo ChainRepository, rpcConfig map[string]*config.RPCDetails) *Chain {
	return &Chain{
		chainRepo: repo,
		rpcConfig: rpcConfig,
	}
}

func (s *Chain) VerifyChain(ctx context.Context, chain string) (bool, error) {
	_, err := s.chainRepo.Exist(ctx, chain)
	if err != nil {
		return false, err
	}

	if chainDetails, ok := s.rpcConfig[chain]; ok {
		if !chainDetails.Enabled {
			return false, err
		}
	} else {
		return false, err
	}

	return true, nil
}

func (s *Chain) ChainName(ctx context.Context, chainId *big.Int) (string, error) {
	name, err := s.chainRepo.GetNameByChainId(ctx, chainId)
	if err != nil {
		return "", err
	}

	return name, nil
}
