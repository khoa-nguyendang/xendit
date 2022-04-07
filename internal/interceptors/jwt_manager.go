package interceptors

import (
	"errors"
	"fmt"
	"time"
	"xendit/config"
	ent "xendit/entities"
	"xendit/pkg/logger"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

type JWTManager struct {
	logger logger.Logger
	cfg    *config.Config
}

func NewJWTManager(logger logger.Logger, cfg *config.Config) *JWTManager {
	return &JWTManager{logger: logger, cfg: cfg}
}

func (manager *JWTManager) GenerateTokenForUser(user *ent.User) (string, error) {

	if user == nil {
		return "", errors.New("staff nil")
	}

	claims := ent.UserClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().UTC().Add(time.Duration(manager.cfg.Server.JwtExpireInHour * int(time.Hour))).Unix(),
		},
		Id:     user.Id,
		RoleId: user.RoleId,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(manager.cfg.Server.JwtSecretKey))
}

// get refresh token.
//
// @params "uuid": device uuid or staff uuid.
//
// @params "expiry": unix milisecond of expiry.
func (manager *JWTManager) GenerateRefreshToken(udid string) (string, error) {
	if udid == "" {
		return "", errors.New("uuid is required")
	}

	phase := fmt.Sprintf("%v:%v", udid, uuid.NewString())
	manager.logger.Infof("generated refresh token phase: %s", phase)
	return phase, nil
}

func (manager *JWTManager) VerifyUser(accessToken string) (*ent.UserClaims, error) {
	token, err := jwt.ParseWithClaims(
		accessToken,
		&ent.UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, fmt.Errorf("unexpected token signing method")
			}
			return []byte(manager.cfg.Server.JwtSecretKey), nil
		},
	)

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*ent.UserClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}
