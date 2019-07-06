package travel

type TripInfo struct {
	TripId    string `json:"tripId"`
	ProfileId string `json:"profileId"`
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
	StepId            int      `json:"stepId"`
	StepName          string   `json:"stepName"`
	ImageBase64       []string `json:"imageBase64"`
	MasterImageBase64 string   `json:"masterImageBase64"`
	MasterMarker      Region   `json:"masterMarker"`
	Markers           []Region `json:"markers"`
	MeanLatitude      float64  `json:"meanLatitude"`
	MeanLongitude     float64  `json:"meanLongitude"`
	Location          string   `json:"location"`
	StartTimestamp    string   `json:"startTimestamp"`
	EndTimestamp      string   `json:"endTimestamp"`
	TimelineData      []string `json:"timelineData"`
	DistanceTravelled int      `json:"distanceTravelled"`
	Description       string   `json:"description"`
	Temperature       string   `json:"temperature"`
}

// Don't forget to update TripUpdateFilter in SaveTrip API
type Trip struct {
	TripId            string   `json:"tripId"`
	TripName          string   `json:"tripName"`
	ProfileId         string   `json:"profileId"`
	Steps             []Step   `json:"steps"`
	Public            bool     `json:"public"`
	MasterPicBase64   string   `json:"masterPicBase64"`
	StartDate         string   `json:"startDate"`
	EndDate           string   `json:"endDate"`
	Temperature       string   `json:"temperature"`
	CountryCode       []string `json:"countryCode"`
	DaysOfTravel      int      `json:"daysOfTravel"`
	DistanceTravelled int      `json:"distanceTravelled"`
	Activities        []string `json:"activities"`
	Location          Region   `json:"location"`
	SyncComplete      bool     `json:"syncComplete"`
}

type TripMeta struct {
	TripId            string   `json:"tripId"`
	TripName          string   `json:"tripName"`
	ProfileId         string   `json:"profileId"`
	Public            bool     `json:"public"`
	MasterPicBase64   string   `json:"masterPicBase64"`
	StartDate         string   `json:"startDate"`
	EndDate           string   `json:"endDate"`
	Temperature       string   `json:"temperature"`
	CountryCode       []string `json:"countryCode"`
	DaysOfTravel      int      `json:"daysOfTravel"`
	DistanceTravelled int      `json:"distanceTravelled"`
	Activities        []string `json:"activities"`
	Location          Region   `json:"location"`
	SyncComplete      bool     `json:"syncComplete"`
}
