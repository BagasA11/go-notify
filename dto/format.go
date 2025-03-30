package dto

type Body struct {
	Receiver string `json:"receiver" binding:"required"`
	Message  string `json:"message"`
}

type Response struct {
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}
