package chat

type ChatMessage struct {
	ProfileId  string `json:"profileid"`
	ChatroomID string `json:"chatroomId"`
	Timestamp  string `json:"timestamp"`
	Message    string `json:"message"`
}

type ChatMessageResponse struct {
	ProfileName string `json:"profileName"`
	ChatroomID  string `json:"chatroomId"`
	Timestamp   string `json:"timestamp"`
	Message     string `json:"message"`
}

type ConnectionMessage struct {
	Timestamp  string `json:"timestamp"`
	ChatroomID string `json:"chatroomId"`
}
