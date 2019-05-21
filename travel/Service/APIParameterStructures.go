package profile

type TravelInfo struct {
	ProfileId int64 `json:"profileId"`
	TravelId  int64 `json:"travelId"`
}

type Region struct {
	Latitude       float64 `json:"latitude"`
	Longitude      float64 `json:"longitude"`
	LatitudeDelta  float64 `json:"latitudeDelta"`
	LongitudeDelta float64 `json:"longitudeDelta"`
}

type ImageData struct {
	Location  Region `json:"region"`
	Image     string `json:"image"`
	Timestamp int64  `json:"timestamp"`
}

type TravelResponse struct {
	Region       Region           `json:"region"`
	ImageData    []ImageData      `json:"imageData"`
	TimelineData map[int][]string `json:"timelinedata"`
}
