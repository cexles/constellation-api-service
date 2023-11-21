package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"math/big"
)

type Chain struct {
	conn *pgxpool.Pool
}

func NewChain(conn *pgxpool.Pool) *Chain {
	return &Chain{
		conn: conn,
	}
}

func (r *Chain) Exist(ctx context.Context, name string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM chains WHERE name = $1);`

	err := r.conn.QueryRow(ctx, query, name).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *Chain) GetNameByChainId(ctx context.Context, chainId *big.Int) (string, error) {
	var name string
	query := `SELECT name FROM chains WHERE chain_id = $1;`

	row := r.conn.QueryRow(ctx, query, chainId)

	err := row.Scan(&name)
	if err != nil {
		return "", err
	}

	return name, nil
}
