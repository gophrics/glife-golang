package travel

type TripInfo struct {
	TripId int `json:"tripId"`
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

type Step struct {
	StepId  int         `json:"stepId"`
	Images  []ImageData `json:"imageData"`
	LatLong Region      `json:"latlong"`
}

type Trip struct {
	TripId    int    `json:"tripId"`
	ProfileId string `json:"profileId"`
	Steps     []Step `json:"steps"`
}
