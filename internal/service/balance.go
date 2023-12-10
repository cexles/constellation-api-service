package service

import (
	"api-service/internal/model"
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/rs/zerolog/log"
	"math/big"
	"sync"
)

type Balance struct {
	ethClients   map[string]*ethclient.Client
	chainService *Chain
	tokenService *Token
}

func NewBalance(ethClients map[string]*ethclient.Client, chainService *Chain, tokenService *Token) *Balance {
	return &Balance{
		ethClients:   ethClients,
		chainService: chainService,
		tokenService: tokenService,
	}
}

func (s *Balance) CoinBalance(ctx context.Context, address string) ([]*model.Balance, error) {
	var (
		balancesMutex sync.Mutex
		balances      []*model.Balance
		wg            sync.WaitGroup
	)

	for _, client := range s.ethClients {
		wg.Add(1)
		go func(client *ethclient.Client) {
			defer wg.Done()

			chainId, err := client.ChainID(ctx)
			if err != nil {
				log.Warn().Msg("Can't parse chain id")
				return
			}

			chainName, err := s.chainService.ChainName(ctx, chainId)
			if err != nil {
				log.Warn().Msg("Can't parse chain name")
				return
			}

			coin, err := s.tokenService.Coin(ctx, chainId)
			if err != nil {
				log.Warn().Msg("Can't parse coin")
				return
			}

			var result string
			err = client.Client().CallContext(ctx, &result, "eth_getBalance", address, "latest")
			if err != nil {
				log.Error().Err(err).Str("address", address).Msg("Failed to get eth balance")
				return
			}

			if len(result) < 2 {
				log.Warn().Msg("Received unexpectedly short balance string")
				return
			}

			balanceBigInt, ok := new(big.Int).SetString(result[2:], 16)
			if !ok || balanceBigInt.Cmp(big.NewInt(0)) == 0 {
				log.Warn().Msg("Can't parse eth balance or balance is zero")
				return
			}

			balancesMutex.Lock()
			balances = append(balances, &model.Balance{
				Chain:    chainName,
				Balance:  balanceBigInt,
				Name:     coin.Name,
				Symbol:   coin.Symbol,
				Decimals: coin.Decimals,
			})
			balancesMutex.Unlock()
		}(client)
	}

	wg.Wait()

	return balances, nil
}

func (s *Balance) TokenBalance(ctx context.Context, address string) ([]*model.Balance, error) {
	var (
		balancesMutex sync.Mutex
		balances      []*model.Balance
		wg            sync.WaitGroup
	)

	encodedCall := EncodeBalanceOfCall(address)

	for _, client := range s.ethClients {
		wg.Add(1)
		go func(client *ethclient.Client) {
			defer wg.Done()

			chainId, err := client.ChainID(ctx)
			if err != nil {
				log.Warn().Msg("Can't parse chain id")
				return
			}

			chainName, err := s.chainService.ChainName(ctx, chainId)
			if err != nil {
				log.Warn().Msg("Can't parse chain name")
				return
			}

			tokens, err := s.tokenService.Tokens(ctx, chainId)
			if err != nil {
				log.Warn().Msg("Can't parse tokens")
				return
			}

			if len(tokens) == 0 {
				return
			}

			batchElems := make([]rpc.BatchElem, len(tokens))

			for i, token := range tokens {
				callArgs := map[string]interface{}{
					"to":   token.Address.Hex(),
					"data": encodedCall,
				}
				batchElems[i] = rpc.BatchElem{
					Method: "eth_call",
					Args:   []interface{}{callArgs, "latest"},
					Result: new(string),
				}
			}

			if err := client.Client().BatchCallContext(ctx, batchElems); err != nil {
				log.Warn().Msgf("BatchCallContext failed: %v", err)
				return
			}

			for i, elem := range batchElems {
				if elem.Error != nil {
					continue
				}

				var result string

				if ptrResult, isPointer := elem.Result.(*string); isPointer {
					result = *ptrResult
				}

				balanceBigInt, ok := new(big.Int).SetString(result[2:], 16)
				if !ok {
					log.Warn().Msg("Can't parse balanceBigInt")
					continue
				}

				if balanceBigInt.Cmp(big.NewInt(0)) == 0 {
					continue
				}

				balancesMutex.Lock()
				balances = append(balances, &model.Balance{
					Chain:    chainName,
					Balance:  balanceBigInt,
					Name:     tokens[i].Name,
					Symbol:   tokens[i].Symbol,
					Decimals: tokens[i].Decimals,
				})
				balancesMutex.Unlock()
			}
		}(client)
	}

	wg.Wait()

	return balances, nil
}

func EncodeBalanceOfCall(address string) string {
	functionSignature := "balanceOf(address)"
	hash := crypto.Keccak256([]byte(functionSignature))
	functionID := hash[:4]

	paddedAddress := common.LeftPadBytes(common.HexToAddress(address).Bytes(), 32)

	data := append(functionID, paddedAddress...)
	return fmt.Sprintf("0x%s", common.Bytes2Hex(data))
}
