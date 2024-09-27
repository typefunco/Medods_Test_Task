package utils

import (
	"context"
	"medods_auth/authService/internal/entities"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(connectionString string) (*PostgresRepository, error) {
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, connectionString)
	if err != nil {
		return nil, err
	}
	return &PostgresRepository{pool: pool}, nil
}

func (r *PostgresRepository) GetUserByID(id int) (*entities.User, error) {
	user := &entities.User{}
	err := r.pool.QueryRow(context.Background(), "SELECT id, username, email FROM users WHERE id = $1", id).
		Scan(&user.ID, &user.Username, &user.Email)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // Пользователь не найден
		}
		return nil, err
	}
	return user, nil
}

func (r *PostgresRepository) SaveRefreshToken(token entities.RefreshToken) error {
	_, err := r.pool.Exec(context.Background(),
		"INSERT INTO refresh_tokens (token_hash, user_id, ip, exp) VALUES ($1, $2, $3, $4)",
		token.TokenHash, token.UserID, token.IP, token.Exp)
	return err
}

func (r *PostgresRepository) GetRefreshTokenByToken(tokenHash string) (*entities.RefreshToken, error) {
	token := &entities.RefreshToken{}
	err := r.pool.QueryRow(context.Background(),
		"SELECT token_hash, user_id, ip, exp FROM refresh_tokens WHERE token_hash = $1", tokenHash).
		Scan(&token.TokenHash, &token.UserID, &token.IP, &token.Exp)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // Токен не найден
		}
		return nil, err
	}
	return token, nil
}

func (r *PostgresRepository) DeleteRefreshToken(tokenHash string) error {
	_, err := r.pool.Exec(context.Background(), "DELETE FROM refresh_tokens WHERE token_hash = $1", tokenHash)
	return err
}

func (r *PostgresRepository) CleanupExpiredTokens() error {
	_, err := r.pool.Exec(context.Background(), "DELETE FROM refresh_tokens WHERE exp < $1", time.Now())
	return err
}
