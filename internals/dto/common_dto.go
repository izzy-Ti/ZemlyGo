package dto

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func SuccessResponse(data interface{}, message string) APIResponse {
	return APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
}

func ErrorAPIResponse(err string) APIResponse {
	return APIResponse{
		Success: false,
		Error:   err,
	}
}

type PaginationMeta struct {
	Limit  int   `json:"limit"`
	Offset int   `json:"offset"`
	Total  int64 `json:"total,omitempty"`
}

type PaginatedResponse struct {
	Items interface{}    `json:"items"`
	Meta  PaginationMeta `json:"meta"`
}
