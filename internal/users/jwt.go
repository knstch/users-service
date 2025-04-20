package users

import (
	"crypto/sha3"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/knstch/subtrack-libs/auth"

	"users-service/internal/domain/enum"
)

func (svc *ServiceImpl) mintJWT(userID uint, role enum.Role) (string, string, error) {
	timeNow := time.Now()
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, auth.Claims{
		UserID: strconv.Itoa(int(userID)),
		Role:   role.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(timeNow.Add(30 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(timeNow),
		},
	})

	signedAccessToken, err := accessToken.SignedString([]byte(svc.jwtSecret))
	if err != nil {
		return "", "", err
	}

	rawRefreshToken := []byte(fmt.Sprintf("%s%d", signedAccessToken, time.Now().Unix()))
	hash := sha3.New256()
	_, err = hash.Write(rawRefreshToken)
	if err != nil {
		return "", "", err
	}
	refreshToken := hex.EncodeToString(hash.Sum(nil))

	return signedAccessToken, refreshToken, nil
}
