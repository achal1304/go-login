package data

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"time"
)

type Token struct {
	UserId    int
	Email     string
	Expiry    time.Time
	Hash      []byte
	Plaintext string
}

func GenerateToken(userID int, ttl time.Duration, email string) (*Token, error) {
	token := &Token{
		UserId: userID,
		Expiry: time.Now().Add(ttl),
		Email:  email,
	}

	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}

	token.Plaintext = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)
	hash := sha256.Sum256([]byte(token.Plaintext))
	token.Hash = hash[:]
	return token, nil
}
