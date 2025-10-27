package handlers

// APIResponse represents a standard API response structure
type APIResponse struct {
	Error    bool        `json:"error"`
	Message  string      `json:"message"`
	Response interface{} `json:"response,omitempty"`
}
