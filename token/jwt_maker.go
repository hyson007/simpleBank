package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const minSecretKeySize = 32

//JWTMaker is a JSON Web Token maker
type JWTMaker struct {
	secretKey string
}

func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size: must be at least %d char", minSecretKeySize)
	}
	return &JWTMaker{secretKey: secretKey}, nil
}

func (j *JWTMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", nil
	}
	// else we create new jwtToken
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return jwtToken.SignedString([]byte(j.secretKey))

}

func (j *JWTMaker) VerifyToken(token string) (*PayLoad, error) {
	//keyfunc receive the parsed but unverified token, we should apply
	//some logic there
	keyFunc := func(jt *jwt.Token) (interface{}, error) {
		// try to convert to specific implementation
		_, ok := jt.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			// means the algorithem doesn't match with our signing method
			return nil, ErrInvalidToken
		}
		return []byte(j.secretKey), nil
	}
	jwtToken, err := jwt.ParseWithClaims(token, &PayLoad{}, keyFunc)

	if err != nil {
		// either token expired or invalid, need to differentiate between them
		// jwt hide our validor error to its inner struct
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}
	//if no error
	payload, ok := jwtToken.Claims.(*PayLoad)
	if !ok {
		return nil, ErrInvalidToken
	}

	return payload, nil
}
