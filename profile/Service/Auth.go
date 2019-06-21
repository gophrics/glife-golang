package profile

import (
	"fmt"
	"net/http"

	"../../common"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
)

func RefreshToken(w http.ResponseWriter, r *http.Request) {
	_, claims, err2 := jwtauth.FromContext(r.Context())

	if err2 != nil {
		fmt.Printf("%s", err2.Error())
		return
	}

	// If struct not initialzed, inner variables don't exist
	profileID := claims["profileid"]

	newClaim := jwt.MapClaims{"profileid": profileID}

	_, token, err := tokenAuth.Encode(newClaim)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	var response common.Token
	response.Token = token
	render.JSON(w, r, response)
}
