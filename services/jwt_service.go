package services

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)


type (
    JWTService interface {
        GenerateAccessToken(userId string, role string) (string, error) 
        GenerateRefreshToken() (string, error)
        ValidateToken(token string) (*jwt.Token, error)
    }

    jwtService struct {
        secretKey []byte
        issuer string
    }
)

func NewJWTService() JWTService {
    return &jwtService{
        secretKey: []byte(os.Getenv("AUTH_SECRET")),
        issuer: os.Getenv("ISSUER"),
    }
}


func (j *jwtService) GenerateAccessToken(userId string, role string) (string, error) {
    // Create access claims
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims {
        "sub": userId,
        "exp": jwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
        "iat": jwt.NewNumericDate(time.Now()),
        "iss": "this-backend-needs-an-iss-id",
        "role": "user-role",
        "type": "access",
    })
    fmt.Println("Access Token Generation", token)

    // Sign the access token
    tokenString, err := token.SignedString(j.secretKey)
    if err != nil {
        return "", err
    }
    fmt.Println("Signed Access Token Generation", tokenString)
    return tokenString, nil
}

func (j *jwtService) GenerateRefreshToken() (string, error) {
    // Create refresh claims (last one day)
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims {
        "exp": jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
        "iat": jwt.NewNumericDate(time.Now()),
        "iss": "this-backend-needs-an-iss-id",
        "type": "refresh",
    })

    // Sign the refresh token
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
