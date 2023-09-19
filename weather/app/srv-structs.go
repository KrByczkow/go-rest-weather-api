package app

type ErrorMessage struct {
	ErrorCode    int    `json:"code"`
	ErrorMessage string `json:"message"`
}
