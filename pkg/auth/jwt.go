package auth

import (
	"fmt"
	"time"

	"github.com/Zhaoyikaiii/docmind/pkg/config"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

type TokenInfo struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

// GenerateToken 生成JWT token
func GenerateToken(userID uint, username, role string) (*TokenInfo, error) {
	secretKey := []byte(config.GetString("jwt.secret"))
	accessExpiration := config.GetInt64("jwt.access_expiration")   // 默认15分钟
	refreshExpiration := config.GetInt64("jwt.refresh_expiration") // 默认7天

	// 创建 access token
	accessClaims := &Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(accessExpiration) * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "docmind",
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString(secretKey)
	if err != nil {
		return nil, fmt.Errorf("could not generate access token: %w", err)
	}

	refreshClaims := &Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(refreshExpiration) * time.Hour * 24)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "docmind",
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString(secretKey)
	if err != nil {
		return nil, fmt.Errorf("could not generate refresh token: %w", err)
	}

	return &TokenInfo{
		AccessToken:  accessTokenString,
		TokenType:    "Bearer",
		ExpiresIn:    accessExpiration * 60,
		RefreshToken: refreshTokenString,
	}, nil
}

func ValidateToken(tokenString string) (*Claims, error) {
	secretKey := []byte(config.GetString("jwt.secret"))

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("could not parse token: %w", err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func RefreshToken(refreshToken string) (*TokenInfo, error) {
	claims, err := ValidateToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	return GenerateToken(claims.UserID, claims.Username, claims.Role)
}
