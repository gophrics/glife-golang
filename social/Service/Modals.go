package social

type FollowRequest struct {
	Following []string `json:"following"`
}

type LikeRequest struct {
	TripId     string `json:"tripId"`
	StepId     string `json:"stepId"`
	ProfileId  string `json:"profileId"`
	TripOrStep int    `json:"triporstep"` // 0 - trip, 1 - step
}

type Step struct {
	StepId    string   `json:"stepId"`
	TripId    string   `json:"tripId"`
	ProfileId string   `json:"profileId"`
	LikedBy   []string `json:"likedBy"`
}

type Trip struct {
	TripId    string   `json:"tripId"`
	ProfileId string   `json:"profileId"`
	LikedBy   []string `json:"likedBy"`
}
