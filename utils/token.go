package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"time"

	"shuttle/databases"
	"shuttle/logger"
	"shuttle/models/entity"
	"shuttle/repositories"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
)

var jwtSecret []byte
var encryptionKey []byte
var db *sqlx.DB

func init() {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	jwtSecret = []byte(viper.GetString("JWT_SECRET"))
	encryptionKey = []byte(viper.GetString("ENCRYPTION_KEY"))

	db, err = databases.PostgresConnection()
	if err != nil {
		panic(err)
	}
}

// Signed Access Token
func GenerateToken(userID, userUUID, username, role_code string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":       userID,
		"user_uuid": userUUID,
		"user_name": username,
		"role_code": role_code,
		"exp":       time.Now().Add(time.Hour * 6).Unix(), // 2 hours expiration
	})

	signedToken, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	encryptedToken, err := encryptToken(signedToken)
	if err != nil {
		return "", err
	}

	return encryptedToken, nil
}

// Same, but with 15 days expiration time and for reissuing access token
func GenerateRefreshToken(userID, userUUID, username, role_code string) (string, error) {

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":       userID,
		"user_uuid": userUUID,
		"user_name": username,
		"role_code": role_code,
		"exp":       time.Now().Add(time.Hour * 24 * 15).Unix(), // 15 days expiration
	})

	signedRefreshToken, err := refreshToken.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	encryptedRefreshToken, err := encryptToken(signedRefreshToken)
	if err != nil {
		return "", err
	}

	return encryptedRefreshToken, nil
}

// AES encryption for tokens
func encryptToken(token string) (string, error) {
	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	encryptedToken := gcm.Seal(nonce, nonce, []byte(token), nil)
	return base64.URLEncoding.EncodeToString(encryptedToken), nil
}

func decryptToken(encryptedToken string) (string, error) {
	encryptedBytes, err := base64.URLEncoding.DecodeString(encryptedToken)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(encryptedBytes) < nonceSize {
		return "", errors.New("malformed encrypted token")
	}

	nonce, ciphertext := encryptedBytes[:nonceSize], encryptedBytes[nonceSize:]
	decryptedToken, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(decryptedToken), nil
}

func ValidateToken(encryptedToken string) (jwt.MapClaims, error) {
	decryptedToken, err := decryptToken(encryptedToken)
	if err != nil {
		return nil, err
	}

	token, err := jwt.Parse(decryptedToken, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}

func SaveRefreshToken(userUUID string, refreshToken string) error {
	ID := time.Now().UnixMilli()*1e6 + int64(uuid.New().ID()%1e6)
	expiration := time.Now().Add(time.Hour * 24 * 15)

	parsedUUID, parseErr := uuid.Parse(userUUID)
	if parseErr != nil {
		return parseErr
	}

	err := repositories.SaveRefreshToken(*db, entity.RefreshToken{
		ID:           ID,
		UserUUID:     parsedUUID,
		RefreshToken: refreshToken,
		ExpiredAt:    expiration,
	})
	if err != nil {
		logger.LogError(err, "Failed to save refresh token", map[string]interface{}{
			"user_id": userUUID,
		})
		return err
	}

	return nil
}

var InvalidTokens = make(map[string]struct{})

func InvalidateToken(token string) {
	const bearerPrefix = "Bearer "
	if len(token) > len(bearerPrefix) && token[:len(bearerPrefix)] == bearerPrefix {
		token = token[len(bearerPrefix):]
	}
	InvalidTokens[token] = struct{}{}
}
