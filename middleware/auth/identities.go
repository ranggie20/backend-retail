package auth

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtSecret string = "ch@ng3Th!$53CRe7"
var jwtTTL time.Duration = 5 * time.Minute
var jwtRefreshTTL time.Duration = 24 * time.Hour

func SetJWTConfig(secret string, ttl time.Duration, refreshTTL time.Duration) {
	if secret != "" {
		jwtSecret = secret
	}
	if ttl > 0 {
		jwtTTL = ttl
	}
	if refreshTTL > 0 {
		jwtRefreshTTL = refreshTTL
	}
	_ = jwtRefreshTTL // TODO: Remove these 2 lines when implemented
}

type (
	jwtClaim struct {
		jwt.StandardClaims
		Data *Identity
	}

	Identity struct {
		UserID       string  `json:"user_id"`
		ProviderName string  `json:"providerName"`
		ProviderType string  `json:"providerType"`
		Issuer       *string `json:"issuer"`
		Primary      string  `json:"primary"`
		DateCreated  string  `json:"dateCreated"`
		Email        string  `json:"email"`
	}
)

// DecodeJWT decodes a JWT token and returns the claims
// func DecodeJWT(tokenString string) (claim Identity, err error) {
// 	var keyFn jwt.Keyfunc = func(t *jwt.Token) (interface{}, error) {
// 		return []byte(jwtSecret), nil
// 	}

// 	var claims jwtClaim
// 	claims.Data = &claim
// 	_, err = jwt.ParseWithClaims(tokenString, &claims, keyFn)
// 	if err == nil {
// 		return
// 	}

// 	err = claims.Valid()
// 	if err != nil {
// 		return
// 	}

//		return
//	}
func DecodeJWT(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, errors.New("invalid token signature")
		}
		return nil, errors.New("invalid token")
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}

// func EncodeJWT(userID, email string) (string, error) {
// 	now := time.Now()

// 	identity := Identity{
// 		UserID:      userID,
// 		Email:       email,
// 		DateCreated: time.Now().Format(time.RFC3339),
// 	}

// 	claims := jwtClaim{
// 		StandardClaims: jwt.StandardClaims{
// 			NotBefore: now.Unix(),
// 			IssuedAt:  now.Unix(),
// 			ExpiresAt: now.Add(jwtTTL).Unix(),
// 		},
// 		Data: &identity,
// 	}

// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	return token.SignedString([]byte(jwtSecret))
// }

func EncodeJWT(user_id int32, username string) (string, error) {
	expirationTime := time.Now().Add(15 * time.Minute) // Set token expiration time
	claims := &Claims{
		UserID:   user_id,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}
