package profile

type User struct {
	ProfileId int64  `json:"profileId"`
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	Country   string `json:"country"`
}
