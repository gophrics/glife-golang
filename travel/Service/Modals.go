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
	StepId        int      `json:"stepId"`
	StepName      string   `json:"stepName"`
	Images        []string `json:"_imageBase64"`
	MasterImage   string   `json:"_masterImageBase64"`
	MasterMarker  Region   `json:"masterMarker"`
	Markers       []Region `json:"markers"`
	MeanLatitude  float64  `json:"meanLatitude"`
	MeanLongitude float64  `json:"meanLongitude"`
}

type Trip struct {
	TripId      int    `json:"tripId"`
	TripName    string `json:"tripName"`
	ProfileId   string `json:"profileId"`
	Steps       []Step `json:"steps"`
	Public      bool   `json:"public"`
	MasterImage string `json:"_masterPicBase64"`
}
