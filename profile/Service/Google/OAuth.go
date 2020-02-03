package profile

import (
	"fmt"
	"net/http"
)

type OAuthHelper interface {
	GetDetails(string) bool
}

// Google asd
type Google struct {
}

func (g Google) VerifyToken(token string) bool {
	_, err := http.Get(fmt.Sprintf("https://oauth2.googleapis.com/tokeninfo?id_token=%s", token))
	return err != nil
}

func (g Google) GetDetails(token string) bool {
	res, err := http.Get(fmt.Sprintf("https://oauth2.googleapis.com/tokeninfo?id_token=%s", token))

}
