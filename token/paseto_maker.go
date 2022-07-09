package token

import (
	"fmt"
	"time"

	"github.com/o1egl/paseto"
	"golang.org/x/crypto/chacha20poly1305"
)

type PasetoMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

func NewPasetoMaker(symkey string) (Maker, error) {
	if len(symkey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size: must be exactly %d chars", chacha20poly1305.KeySize)
	}

	return &PasetoMaker{paseto: paseto.NewV2(), symmetricKey: []byte(symkey)}, nil
}

//CreateToken creates a new token for a specific username and duration
func (p *PasetoMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}

	return p.paseto.Encrypt(p.symmetricKey, payload, nil)
}

//VerifyToken checks if the token is valid or not
func (p *PasetoMaker) VerifyToken(token string) (*PayLoad, error) {
	payload := &PayLoad{}
	err := p.paseto.Decrypt(token, p.symmetricKey, payload, nil)
	if err != nil {
		return nil, ErrInvalidToken
	}
	err = payload.Valid()
	if err != nil {
		return nil, err
	}
	return payload, nil
}
