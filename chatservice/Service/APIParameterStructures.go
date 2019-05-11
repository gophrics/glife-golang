package location

type ChatMessage struct {
	ProfileId  string `json:"profileId"`
	ChatroomID string `json:"chatroomId"`
	Timestamp  string `json:"timestamp"`
	Message    string `json:"message"`
}
