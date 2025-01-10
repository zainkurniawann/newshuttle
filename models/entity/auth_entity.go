package entity

import (
	"time"

	"github.com/google/uuid"
)

type UserDataOnLogin struct {
	ID       int64  `db:"user_id"`
	UUID     string `db:"user_uuid"`
	Username string `db:"user_username"`
	RoleCode string `db:"user_role_code"`
	Password string `db:"user_password"`
}

type RefreshToken struct {
	ID           int64     `db:"id"`
	UserUUID     uuid.UUID `db:"user_uuid"`
	RefreshToken string    `db:"refresh_token"`
	IssuedAt     time.Time `db:"issued_at"`
	ExpiredAt    time.Time `db:"expired_at"`
	Revoked      bool      `db:"is_revoked"`
	LastUsedAt   *time.Time `db:"last_used_at"`
}

type FCMToken struct {
	ID          int64     `db:"id"`
	UserUUID    uuid.UUID `db:"user_uuid"`
	DeviceToken string    `db:"device_token"`
	CreatedAt   time.Time `db:"created_at"`
	UpdateAt    time.Time `db:"updated_at"`
}
