package profile

import (
	"context"
	"fmt"
	"net/http"

	"../../common/mongodb"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var tokenAuth *jwtauth.JWTAuth

func init() {
	tokenAuth = jwtauth.New("HS256", []byte("dasjkhfiuadufasdfasf832742389r923rc325235c7n6235"), nil)

	// For debugging/example purposes, we generate and print
	// a sample jwt token with claims `user_id:123` here:
	_, tokenString, _ := tokenAuth.Encode(jwt.MapClaims{"user_id": 123})
	fmt.Printf("DEBUG: a sample jwt is %s\n\n", tokenString)
}

func Routes() *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.Recoverer)

	// Protected routes
	router.Group(func(r chi.Router) {
		// Seek, verify and validate JWT tokens
		r.Use(jwtauth.Verifier(tokenAuth))

		// Handle valid / invalid tokens. In this example, we use
		// the provided authenticator middleware, but you can write your
		// own very easily, look at the Authenticator method in jwtauth.go
		// and tweak it, its not scary.
		r.Use(jwtauth.Authenticator)
		r.Get("/api/v1/profile/search/{searchstring}", FindUser)
		r.Get("/api/v1/profile/getuser", GetUser)
	})

	// Public routes
	router.Group(func(r chi.Router) {
		r.Post("/api/v1/profile/register", RegisterUser)
		r.Post("/api/v1/profile/login", LoginUser)
		r.Post("/api/v1/profile/registerWithGoogle", RegisterUserWithGoogle)
		r.Post("/api/v1/profile/loginWithGoogle", LoginUserWithGoogle)
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("welcome anonymous"))
		})
	})

	return router
}

func GetUser(w http.ResponseWriter, r *http.Request) {

	token, claims, err2 := jwtauth.FromContext(r.Context())

	if err2 != nil {
		fmt.Printf("%s", err2.Error())
		return
	}

	fmt.Printf("%s", token)
	fmt.Printf("%s", claims)
	// If struct not initialzed, inner variables don't exist
	profileId := fmt.Sprintf("%s", "nitin.i.joy@gmail.com")

	filter := bson.D{{
		"_id",
		bson.D{{
			"$in",
			bson.A{profileId},
		}},
	}}

	var result GetUserResponse

	err := mongodb.Profile.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusGone)
	}

	render.JSON(w, r, result)
}

func FindUser(w http.ResponseWriter, r *http.Request) {

	fmt.Sprintf("FindUser called")
	ss := fmt.Sprintf("%s", chi.URLParam(r, "searchstring"))
	ss = fmt.Sprintf(".*%s.*", ss)

	filter := bson.D{
		{"$or", bson.A{
			bson.D{{"name", primitive.Regex{ss, ""}}},
			bson.D{{"country", ss}},
			bson.D{{"phone", ss}},
		}}}
	findOptions := options.Find()

	var result []GetUserResponse

	cur, err := mongodb.Profile.Find(context.TODO(), filter, findOptions)

	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	var x GetUserResponse
	for cur.Next(context.TODO()) {
		err := cur.Decode(&x)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		fmt.Printf("%s", x)
		result = append(result, x)
	}

	render.JSON(w, r, result)

}
