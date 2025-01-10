package repositories

import (
	"context"
	"shuttle/models/entity"
	"time"

	"github.com/jmoiron/sqlx"
)

type AuthRepositoryInterface interface {
	Login(email string) (entity.UserDataOnLogin, error)
	CheckRefreshTokenData(userUUID, token string) (entity.RefreshToken, error)
	DeleteRefreshToken(ctx context.Context, userUUID string) error
	UpdateUserStatus(userUUID, status string, lastActive time.Time) error
	UpdateRefreshToken(userUUID, refreshToken string) (time.Time, error)
	SaveDeviceToken(tokendata entity.FCMToken) error
}

type authRepository struct {
	DB *sqlx.DB
}

func NewAuthRepository(DB *sqlx.DB) AuthRepositoryInterface {
	return &authRepository{
		DB: DB,
	}
}

func (r *authRepository) Login(email string) (entity.UserDataOnLogin, error) {
	var user entity.UserDataOnLogin
	query := `
		SELECT
		user_id, user_uuid, user_username, user_role_code, user_password
		FROM users 
		WHERE user_email = $1 AND deleted_at IS NULL
	`

	row := r.DB.QueryRow(query, email)

	if err := row.Scan(&user.ID, &user.UUID, &user.Username, &user.RoleCode, &user.Password); err != nil {
		return entity.UserDataOnLogin{}, err
	}

	return user, nil
}

func (r *authRepository) CheckRefreshTokenData(userUUID, token string) (entity.RefreshToken, error) {
	query := `
		SELECT refresh_token, expired_at, is_revoked, last_used_at
		FROM refresh_tokens 
		WHERE user_uuid = $1 AND refresh_token = $2
	`

	var tokenData entity.RefreshToken
	err := r.DB.Get(&tokenData, query, userUUID, token)
	if err != nil {
		return tokenData, err
	}

	return tokenData, nil
}

func SaveRefreshToken(db sqlx.DB, refreshToken entity.RefreshToken) error {
	query := `
		INSERT INTO refresh_tokens (id, user_uuid, refresh_token, expired_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_uuid)
		DO UPDATE SET refresh_token = $3, issued_at = CURRENT_TIMESTAMP, expired_at = $4, is_revoked = false 
	`
	_, err := db.Exec(query, refreshToken.ID, refreshToken.UserUUID, refreshToken.RefreshToken, refreshToken.ExpiredAt)
	if err != nil {
		return err
	}

	return nil
}

func (r *authRepository) DeleteRefreshToken(ctx context.Context, userUUID string) error {
	query := `
		DELETE FROM refresh_tokens
		WHERE user_uuid = $1
	`

	_, err := r.DB.ExecContext(ctx, query, userUUID)
	if err != nil {
		return err
	}

	return nil
}

func (r *authRepository) UpdateUserStatus(userUUID, status string, lastActive time.Time) error {
	query := `
		UPDATE users
		SET user_status = $1, user_last_active = $2
		WHERE user_uuid = $3
	`

	_, err := r.DB.Exec(query, status, lastActive, userUUID)
	if err != nil {
		return err
	}

	return nil
}

func (r *authRepository) UpdateRefreshToken(userUUID, refreshToken string) (time.Time, error) {
	var lastUsedAt time.Time
	
	querySelect := `
		SELECT last_used_at 
		FROM refresh_tokens 
		WHERE user_uuid = $1 AND refresh_token = $2
	`
	err := r.DB.QueryRow(querySelect, userUUID, refreshToken).Scan(&lastUsedAt)
	if err != nil {
		return time.Time{}, err
	}

	// Query untuk update last_used_at ke NOW()
	queryUpdate := `
		UPDATE refresh_tokens
		SET last_used_at = NOW()
		WHERE user_uuid = $1 AND refresh_token = $2
	`
	_, err = r.DB.Exec(queryUpdate, userUUID, refreshToken)
	if err != nil {
		return time.Time{}, err
	}

	return lastUsedAt, nil
}

func (r *authRepository) SaveDeviceToken(tokendata entity.FCMToken) error {
	query := `
		INSERT INTO fcm_tokens (id, user_uuid, device_token, created_at)
		VALUES ($1, $2, $3, NOW())
		ON CONFLICT (user_uuid)
		DO UPDATE SET device_token = $2, updated_at = NOW()
	`
	_, err := r.DB.Exec(query, tokendata.ID, tokendata.UserUUID, tokendata.DeviceToken)
	if err != nil {
		return err
	}

	return nil
}