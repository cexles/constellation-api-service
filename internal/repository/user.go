package repository

import (
	"api-service/internal/model"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type User struct {
	conn *pgxpool.Pool
}

func NewUser(conn *pgxpool.Pool) *User {
	return &User{
		conn: conn,
	}
}

func (r *User) CreateOrUpdate(ctx context.Context, address string) (model.User, error) {
	u := model.User{
		Address:  address,
		OnlineAt: time.Now().UTC(),
	}
	query := `
		INSERT INTO users (address, online_at) 
		VALUES ($1, $2) 
		ON CONFLICT (address) 
		DO UPDATE SET online_at = $3`

	_, err := r.conn.Exec(ctx, query, u.Address, u.OnlineAt, u.OnlineAt)
	if err != nil {
		return model.User{}, err
	}

	return u, nil
}
