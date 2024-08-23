package services

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)


type (
    JWTService interface {
        GenerateToken(userId string, role string) (string, error) 
        ValidateToken(token string) (*jwt.Token, error)
    }

    jwtService struct {
        secretKey string
        issuer string
    }
)

func NewJWTService() JWTService {
    return &jwtService{
        secretKey: os.Getenv("AUTH_SECRET"),
        issuer: os.Getenv("ISSUER"),
    }
}

func (j *jwtService) GenerateToken(userId string, role string) (string, error) {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user": userId,
        "role": role,
        "exp": time.Now().Add(time.Hour * 1).Unix(),
    })
    tokenString, err := token.SignedString(j.secretKey)
    if err != nil {
        return "", err
    }

    return tokenString, nil
}

func (j *jwtService) ValidateToken (tokenString string) (*jwt.Token, error) {
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }

        return j.secretKey, nil
    })

    if err != nil {
        return nil, err
    }

    if !token.Valid {
        return nil, fmt.Errorf("invalid token")
    }
    return token, nil
}
