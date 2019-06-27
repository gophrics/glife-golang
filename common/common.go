package common

const JWT_SECRET_KEY string = "dasjkhfiuadufasdfasf832742389r923rc325235c7n6235"
const GOOGLE_APP_ID = "249369235819-11cfia1ht584n1kmk6gh6kbba8ab429u.apps.googleusercontent.com"
const LOCATIONIQ_TOKEN = "225053fec17f73"

type Token struct {
	Token string
}

type User struct {
	ProfileId string `json:"profileid"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Country   string `json:"country"`
	Username  string `json:"username"`
	Password  string `json:"password"`
}

type Error struct {
	Type    int    `json:"type"`
	Message string `json:"message"`
}
