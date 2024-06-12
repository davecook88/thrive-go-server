package handlers

type UserMessage struct {
	Message string `json:"message" binding:"required"`
}
