package users

import (
	"crypto/sha256"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID string `json:"uid"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func (svc *ServiceImpl) mintJWT(userID uint, role string) (string, string, error) {
	timeNow := time.Now()
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		UserID: strconv.Itoa(int(userID)),
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(timeNow.Add(30 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(timeNow),
		},
	})

	signedAccessToken, err := accessToken.SignedString([]byte(svc.jwtSecret))
	if err != nil {
		return "", "", err
	}

	refreshToken := string(sha256.New().Sum([]byte(fmt.Sprintf("%s%d", signedAccessToken, time.Now().Unix()))))

	return signedAccessToken, refreshToken, nil
}
