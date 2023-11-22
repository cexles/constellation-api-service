package handler

import (
	"api-service/internal/api/response"
	"api-service/internal/service"
	"github.com/gofiber/fiber/v2"
	"math/big"
	"net/http"
)

type UserApi struct {
	chainService   *service.Chain
	balanceService *service.Balance
}

func NewUserApi(chainService *service.Chain, balanceService *service.Balance) *UserApi {
	return &UserApi{
		chainService:   chainService,
		balanceService: balanceService,
	}
}

func (a *UserApi) Balance(c *fiber.Ctx) error {
	address := c.Query("address")

	coinBalance, err := a.balanceService.CoinBalance(c.Context(), address)
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, "Can't parse coin balance")
	}

	tokenBalance, err := a.balanceService.TokenBalance(c.Context(), address)
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, "Can't parse token balance")
	}

	totalBalance := append(coinBalance, tokenBalance...)

	symbolMap := make(map[string]*response.UserBalance)

	for _, token := range totalBalance {
		balance, exists := symbolMap[token.Symbol]
		if !exists {
			balance = &response.UserBalance{
				Name:         token.Name,
				Symbol:       token.Symbol,
				Decimals:     token.Decimals,
				TokenBalance: new(big.Int),
				Price:        0,
				Chains:       make(map[string]response.ChainInfo),
			}
			symbolMap[token.Symbol] = balance
		}
		balance.TokenBalance.Add(balance.TokenBalance, token.Balance)

		chainInfo, chainExists := balance.Chains[token.Chain]
		if !chainExists {
			chainInfo = response.ChainInfo{
				TokenBalance: new(big.Int),
			}
		}
		chainInfo.TokenBalance.Add(chainInfo.TokenBalance, token.Balance)
		balance.Chains[token.Chain] = chainInfo
	}

	var balances []response.UserBalance
	for _, balance := range symbolMap {
		balances = append(balances, *balance)
	}

	return c.JSON(balances)
}
