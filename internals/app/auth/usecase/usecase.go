package usecase

import (
	"context"
	"crypto/sha1"
	"fmt"
	"time"

	"dclibrary.com/internals/app/auth"
	"dclibrary.com/internals/app/db"
	"dclibrary.com/internals/app/models"
	"github.com/dgrijalva/jwt-go/v4"
)

type Authorizer struct {
	storage *db.UsersStorage

	hashSalt       string
	signingKey     []byte
	expireDuration time.Duration
}

func NewAuthorizer(storage *db.UsersStorage, hashSalt string, signingKey []byte, expireDuration time.Duration) *Authorizer {
	return &Authorizer{
		storage:        storage,
		hashSalt:       hashSalt,
		signingKey:     signingKey,
		expireDuration: expireDuration,
	}
}

func (a *Authorizer) SignUp(user *models.UserAuth) error {
	// Create password hash
	pwd := sha1.New()
	pwd.Write([]byte(user.Password))
	pwd.Write([]byte(a.hashSalt))
	user.Password = fmt.Sprintf("%x", pwd.Sum(nil))

	return a.storage.CreateNewUser(user)
}

func (a *Authorizer) SignIn(ctx context.Context, user *models.UserAuth) (string, error) {
	pwd := sha1.New()
	pwd.Write([]byte(user.Password))
	pwd.Write([]byte(a.hashSalt))
	user.Password = fmt.Sprintf("%x", pwd.Sum(nil))

	user, err := a.repo.Get(ctx, user.Username, user.Password)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &auth.Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.At(time.Now().Add(a.expireDuration)),
			IssuedAt:  jwt.At(time.Now()),
		},
		Username: user.Username,
	})

	return token.SignedString(a.signingKey)
}
