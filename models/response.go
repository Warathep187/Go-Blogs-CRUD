package models

type ErrorResponse struct {
	Message string `json:"message"`
}

type ValidationErrorResponse struct {
	Field   string      `json:"field"`
	Value   interface{} `json:"value"`
	Tag     string      `json:"tag"`
	Message string      `json:"message"`
}

type SuccessResponse struct {
	Message string `json:"message"`
}
