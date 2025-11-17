package auth

import (
	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
	"fmt"
)

func HashPassword(password string) (string, error) {
	// hash the password using argon2id.CreateHash
	hashedpw, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", err
	}
	return hashedpw, nil
}

func CheckPasswordHash(password, hash string) (bool, error) {
	// check the password against the hash
	match, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		return false, err
	}
	return match, nil
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	now := time.Now()
	// define the JWT claims
	claims := &jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(expiresIn)),
		Subject:   userID.String(),
	}
	// create a new token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// sign the token with the secret key
	signedToken, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}
	// return the signed token and nil error
	return signedToken, nil	
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	// ngl i made chatgpt write this entire function it was super confusing
	// but i am writing comments to understand it better
	claims := &jwt.RegisteredClaims{}

	// parse the token with the claims
	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (any, error) {
			return []byte(tokenSecret), nil
		},
	)
	if err != nil {
		return uuid.Nil, err
	}

	// token.Valid checks expiration + signature
	if !token.Valid {
		return uuid.Nil, fmt.Errorf("invalid token")
	}

	// Validate issuer manually if required
	if claims.Issuer != "chirpy" {
		return uuid.Nil, fmt.Errorf("invalid issuer")
	}

	// Extract user ID from Subject
	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, err
	}

	return userID, nil
}