package jwt

import (
	"crypto/rsa"
	"errors"
	"io/ioutil"
	"log"
	"path/filepath"
	"runtime"
	"time"

	"github.com/happilymarrieddad/old-world/api3/types"

	jwtpkg "github.com/dgrijalva/jwt-go"
)

const (
	privKeyPath  = "/../../keys/app.rsa"     // openssl genrsa -out app.rsa keysize
	pubKeyPath   = "/../../keys/app.rsa.pub" // openssl rsa -in app.rsa -pubout > app.rsa.pub
	HOURS_IN_DAY = 24
)

var (
	verifyKey *rsa.PublicKey
	signKey   *rsa.PrivateKey
)

func init() {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	signBytes, err := ioutil.ReadFile(basepath + privKeyPath)
	if err != nil {
		panic(err)
	}
	signKey, err = jwtpkg.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		panic(err)
	}
	verifyBytes, err := ioutil.ReadFile(basepath + pubKeyPath)
	if err != nil {
		panic(err)
	}
	verifyKey, err = jwtpkg.ParseRSAPublicKeyFromPEM(verifyBytes)
	if err != nil {
		panic(err)
	}
}

func NewToken(usr *types.User) string {
	if usr == nil {
		// This should never happen but it's okay because it will just create an invalid user
		usr = new(types.User)
	}

	token := jwtpkg.New(jwtpkg.SigningMethodRS512)
	claims := make(jwtpkg.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * HOURS_IN_DAY).Unix()
	claims["iat"] = time.Now().Unix()
	claims["id"] = usr.ID

	token.Claims = claims

	tokenString, err := token.SignedString(signKey)
	if err != nil {
		panic(err)
	}

	return tokenString
}

func IsTokenValid(val string) (*string, error) {
	token, err := jwtpkg.Parse(val, func(token *jwtpkg.Token) (interface{}, error) {
		return verifyKey, nil
	})

	switch vErr := err.(type) {
	case nil:
		if !token.Valid {
			return nil, errors.New("token is invalid")
		}

		claims, ok := token.Claims.(jwtpkg.MapClaims)
		if !ok {
			return nil, errors.New("token is invalid")
		}

		usrID, ok := claims["id"].(string)
		if !ok {
			return nil, errors.New("unable to get user id from token")
		}

		return &usrID, nil
	case *jwtpkg.ValidationError:
		switch vErr.Errors {
		case jwtpkg.ValidationErrorExpired:
			return nil, errors.New("token expired, get a new one")
		default:
			log.Println(vErr)
			return nil, errors.New("error while parsing token")
		}
	default:
		return nil, errors.New("unable to parse token")
	}
}
