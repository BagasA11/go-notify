package dto

type Body struct {
	Sender   string `json:"sender" binding:"required"`
	Receiver string `json:"receiver" binding:"required"`
	Subject  string `json:"subject" binding:"required"`
	Message  string `json:"message"`
}

type Response struct {
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}
