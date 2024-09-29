package types

type AuthResponse struct {
	Token string `json:"token"`
}

type ErrorResponse struct {
	Timestamp string `json:"timestamp"`
	Message   string `json:"message"`
	RequestID string `json:"requestId"`
}
