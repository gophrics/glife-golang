package location

type Location struct {
	ProfileId string  `json:"profileId"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
