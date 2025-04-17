package response

type SuccessResponse struct {
	Status int `json:"status"`
	Data   any `json:"data,omitempty"`
}

type ErrorResponse struct {
	Status int    `json:"status"`
	Error  string `json:"error"`
}
