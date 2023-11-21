package repository

import (
	"api-service/internal/model"
	"context"
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"math/big"
)

type Token struct {
	conn *pgxpool.Pool
}

func NewToken(conn *pgxpool.Pool) *Token {
	return &Token{
		conn: conn,
	}
}

func (r *Token) GetTokensByChainId(ctx context.Context, chainId *big.Int) ([]model.Token, error) {
	query := `
		SELECT t.address, t.name, t.symbol, t.decimals
		FROM tokens t
		JOIN chains c ON t.chain_id = c.id
		WHERE c.chain_id = $1 AND t.address != '0x0000000000000000000000000000000000000000';
	`
	rows, err := r.conn.Query(ctx, query, chainId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tokens []model.Token
	for rows.Next() {
		var t model.Token
		var address string
		if err := rows.Scan(&address, &t.Name, &t.Symbol, &t.Decimals); err != nil {
			return nil, err
		}
		t.Address = common.HexToAddress(address)
		tokens = append(tokens, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tokens, nil
}

func (r *Token) GetCoinByChainId(ctx context.Context, chainId *big.Int) (model.Token, error) {
	query := `
		SELECT t.address, t.name, t.symbol, t.decimals
		FROM tokens t
		JOIN chains c ON t.chain_id = c.id
		WHERE c.chain_id = $1 AND t.address = '0x0000000000000000000000000000000000000000';
	`
	row := r.conn.QueryRow(ctx, query, chainId)

	var t model.Token
	var address string
	if err := row.Scan(&address, &t.Name, &t.Symbol, &t.Decimals); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Token{}, nil
		}
		return model.Token{}, err
	}

	t.Address = common.HexToAddress(address)
	return t, nil
}
