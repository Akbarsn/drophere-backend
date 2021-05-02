package jwt_tools

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/bccfilkom/drophere-go/domain"
	jwt "github.com/dgrijalva/jwt-go"
)

var (
	// A private key for context that only this package can access. This is important
	// to prevent collisions between different context uses
	userCtxKey      = &contextKey{"user"}
	errInvalidToken = errors.New("jwt: invalid token")
)

type contextKey struct {
	name string
}

// JWTAuthenticator struct
type JWTAuthenticator struct {
	key      []byte
	duration time.Duration
	algo     string
	userRepo domain.UserRepository
}

// NewJWT func
func NewJWT(secret string, duration time.Duration, algo string, userRepo domain.UserRepository) *JWTAuthenticator {
	return &JWTAuthenticator{
		key:      []byte(secret),
		duration: duration,
		algo:     algo,
		userRepo: userRepo,
	}
}

// Authenticate func
func (j *JWTAuthenticator) Authenticate(u *domain.User) (*domain.UserCredentials, error) {
	expiry := time.Now().Add(j.duration)
	token := jwt.NewWithClaims(jwt.GetSigningMethod(j.algo), jwt.MapClaims{
		"user_id": u.ID,
		"exp":     expiry.Unix(),
	})

	tokenS, err := token.SignedString(j.key)
	if err != nil {
		return nil, err
	}
	return &domain.UserCredentials{
		Token:  tokenS,
		Expiry: &expiry,
	}, nil
}

func (j *JWTAuthenticator) validateAndGetUserID(token string) (uint, error) {
	payloadI, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod(j.algo) != token.Method {
			return nil, errInvalidToken
		}

		return j.key, nil
	})

	if err != nil {
		return 0, err
	}

	if !payloadI.Valid {
		return 0, errInvalidToken
	}

	claims := payloadI.Claims.(jwt.MapClaims)

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return 0, errInvalidToken
	}

	return uint(userID), nil
}

func writeGqlError(w http.ResponseWriter, msg string) {
	json.NewEncoder(w).Encode(map[string]interface{}{
		"errors": []map[string]string{
			{
				"message": msg,
			},
		},
	})
}
