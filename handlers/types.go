package handlers


type UserMessage struct {
	ChatID *string `json:"chat_id"`
	Message string `json:"message" binding:"required"`
}