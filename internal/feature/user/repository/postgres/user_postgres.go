package postgres

import (
	"Test/internal/feature/user/entity"
	"Test/internal/feature/user/repository"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type PostgresUserRepository struct {
	db *pgxpool.Pool
}

func NewPostgresUserRepository(db *pgxpool.Pool) repository.UserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) Create(ctx context.Context, user *entity.User) error {
	query := `INSERT INTO users (full_name, username, hashed_password, email, phone, role, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	err := r.db.QueryRow(ctx, query,
		user.FullName,
		user.Username,
		user.HashedPassword,
		user.Email,
		user.Phone,
		user.Role,
		user.CreatedAt,
	).Scan(&user.ID)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (r *PostgresUserRepository) FindByEmail(ctx context.Context, email string) (entity.User, error) {
	query := `SELECT id, full_name, username, hashed_password, email, phone, role, created_at FROM users WHERE email = $1`
	var user entity.User
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.FullName,
		&user.Username,
		&user.HashedPassword,
		&user.Email,
		&user.Phone,
		&user.Role,
		&user.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return entity.User{}, fmt.Errorf("user not found: %w", err)
		}
		return entity.User{}, fmt.Errorf("failed to find user by username: %w", err)
	}
	return user, nil
}

func (r *PostgresUserRepository) FindByPhone(ctx context.Context, phone string) (entity.User, error) {
	query := `SELECT id, full_name, username, hashed_password, email, phone, role, created_at FROM users WHERE phone = $1`
	var user entity.User

	err := r.db.QueryRow(ctx, query, phone).Scan(
		&user.ID,
		&user.FullName,
		&user.Username,
		&user.HashedPassword,
		&user.Email,
		&user.Phone,
		&user.Role,
		&user.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.User{}, err
		}
		return entity.User{}, fmt.Errorf("failed to find user by phone: %w", err)
	}

	return user, nil
}

func (r *PostgresUserRepository) FindById(ctx context.Context, id uint64) (entity.User, error) {
	query := `SELECT id, full_name, username, hashed_password, email, phone, role, created_at FROM users WHERE id = $1`
	var user entity.User
	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.FullName,
		&user.Username,
		&user.HashedPassword,
		&user.Email,
		&user.Phone,
		&user.Role,
		&user.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return entity.User{}, fmt.Errorf("user not found: %w", err)
		}
		return entity.User{}, fmt.Errorf("failed to find user by username: %w", err)
	}
	return user, nil
}

func (r *PostgresUserRepository) FindByUsername(ctx context.Context, username string) (entity.User, error) {
	query := `SELECT id, full_name, username, hashed_password, email, phone, role, created_at FROM users WHERE username = $1`
	var user entity.User
	err := r.db.QueryRow(ctx, query, username).Scan(
		&user.ID,
		&user.FullName,
		&user.Username,
		&user.HashedPassword,
		&user.Email,
		&user.Phone,
		&user.Role,
		&user.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return entity.User{}, fmt.Errorf("user not found: %w", err)
		}
		return entity.User{}, fmt.Errorf("failed to find user by username: %w", err)
	}
	return user, nil
}

func (r *PostgresUserRepository) SaveRefreshToken(ctx context.Context, userID uint64, tokenID string, expiresAt time.Time) error {
	query := `INSERT INTO refresh_tokens (id, user_id, expires_at) VALUES ($1, $2, $3)`
	_, err := r.db.Exec(ctx, query, tokenID, userID, expiresAt)
	if err != nil {
		return fmt.Errorf("failed to save refresh token: %w", err)
	}
	log.Debug().Str("token_id", tokenID).Uint64("user_id", userID).Msg("Refresh token saved")
	return nil
}

func (r *PostgresUserRepository) RevokeRefreshToken(ctx context.Context, tokenID string) error {
	query := `UPDATE refresh_tokens SET revoked_at = NOW() WHERE id = $1`
	res, err := r.db.Exec(ctx, query, tokenID)
	if err != nil {
		return fmt.Errorf("failed to delete refresh token: %w", err)
	}
	rowsAffected := res.RowsAffected()
	if rowsAffected == 0 {
		log.Warn().Str("token_id", tokenID).Msg("Refresh token not found for deletion")
		return fmt.Errorf("refresh token not found: %w", err)
	}
	log.Debug().Str("token_id", tokenID).Msg("Refresh token deleted")
	return nil
}

func (r *PostgresUserRepository) FindRefreshToken(ctx context.Context, tokenID string) (repository.RefreshToken, error) {
	query := `SELECT id, user_id, expires_at, created_at, revoked_at  FROM refresh_tokens WHERE id = $1 AND expires_at > NOW()`
	var dbRefreshToken repository.RefreshToken

	err := r.db.QueryRow(ctx, query, tokenID).Scan(&dbRefreshToken.ID, &dbRefreshToken.UserID, &dbRefreshToken.ExpiresAt, &dbRefreshToken.CreatedAt, &dbRefreshToken.RevokedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return dbRefreshToken, fmt.Errorf("refresh token not found or expired: %w", err)
		}
		return dbRefreshToken, fmt.Errorf("failed to find refresh token: %w", err)
	}
	return dbRefreshToken, nil
}

func (r *PostgresUserRepository) DeleteUser(ctx context.Context, userID uint64) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.Exec(ctx, query, userID)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgresUserRepository) GetOrganizationID(ctx context.Context, userID uint64) (*uint64, error) {
	var organizationID *uint64
	query := `SELECT organization_id FROM organizations_users WHERE user_id = $1`
	err := r.db.QueryRow(ctx, query, userID).Scan(&organizationID)
	if err != nil {
		return nil, err
	}
	return organizationID, nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(bytes), nil
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
