package services

import (
	"os"
	"time"
    "fmt"

	"github.com/golang-jwt/jwt/v5"
)

type (
    JWTService interface {
        GenerateAccessToken(string, string) (string, error) 
        GenerateRefreshToken(string, string) (string, error)
        ValidateToken(string) (*jwt.Token, error)
        RefreshToken(string) (string, error)
        GetTokenClaims(token *jwt.Token) (jwt.MapClaims, error)
    }

    jwtService struct {
        secretKey []byte
        issuer string
    }
)

func NewJWTService() JWTService {
    return &jwtService {
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
        "role": role,
        "type": "access",
    })

    // Sign the access token
    tokenString, err := token.SignedString(j.secretKey)
    if err != nil {
        return "", err
    }
    return tokenString, nil
}

func (j *jwtService) GenerateRefreshToken(userId string, role string) (string, error) {
    // Create refresh claims (last one day)
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims {
        "sub": userId,
        "exp": jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
        "iat": jwt.NewNumericDate(time.Now()),
        "iss": "this-backend-needs-an-iss-id",
        "role": role,
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
            return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
        }

        return j.secretKey, nil
    })

    if err != nil {
        return nil, err
    }

    if !token.Valid {
        return nil, fmt.Errorf("Invalid Token")
    }
    return token, nil
}

func (j *jwtService) GetTokenClaims(token *jwt.Token) (jwt.MapClaims, error) {
    // Make sure the token is valid in the first place
    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok || !token.Valid {
        return nil, fmt.Errorf("Invalid Token")
    }
    return claims, nil
}

// Must be a refresh token when requesting an access token refresh.
// Will help mitigate stupid requests a little bit
func (j *jwtService) RefreshToken(tokenString string) (string, error) {
    token, err := j.ValidateToken(tokenString)
    if err != nil {
        return "", err
    }

    claims, err := j.GetTokenClaims(token)
    if err != nil {
        return "", err
    }

    // Tokens must be refresh tokens to request refresh
    tokenType, ok := claims["type"].(string)
    if !ok || tokenType != "refresh" {
        return "", fmt.Errorf("Invalid token type: Not a refresh token")
    }

    // UserId needed for new access token
    userId, ok := claims["sub"].(string)
    if !ok {
        return "", fmt.Errorf("Invalid Token: User ID is missing")
    }

    // User role needed for new access token
    role, ok := claims["role"].(string)
    if !ok {
        return "", fmt.Errorf("Invalid Token: User Role is missing")
    }

    newAccessToken, err := j.GenerateAccessToken(userId, role)
    if err != nil {
        return "", err
    }

    return newAccessToken, nil
}

