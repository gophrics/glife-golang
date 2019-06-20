package authentication

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("my_secret_key")

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type JWTTokenInfo struct {
	Expires time.Time
	Value   string
}

// A Dumb function that will generate JWT Token for each and every thing
func GenerateJWTToken(username string) (*JWTTokenInfo, error) {
	expirationTime := time.Now().Add(5 * time.Minute)

	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)

	return &JWTTokenInfo{
		Value:   tokenString,
		Expires: expirationTime,
	}, err
}

// A Smart function that will validate if the token is in par with the claims
func VerifyJWTToken(unparsedToken string) (bool, *Claims, error) {

	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(unparsedToken, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if !tkn.Valid {
		return false, &Claims{}, nil
	}

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return false, &Claims{}, nil
		}
		return false, &Claims{}, err
	}

	return true, claims, nil
}
