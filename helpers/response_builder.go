package helpers

// Response creates a JSON object that must contain a message string.
// It can optionally have Error or Data filled with any data.
type Response struct {
	Message string `json:"message"`
	Error   any    `json:"error,omitempty"`
	Data    any    `json:"data,omitempty"`
}

// BuildSuccessfulResponse creates a response type consisting of a message
// string, and any data you want to pass out to the end user.
func BuildSuccessfulResponse(message string, data any) Response {
    return Response {
        Message: message,
        Data: data,
    }
}

// BuildFailedResponse creates a response type consisting of a message
// string, an error of any type, and any data you want to pass out to 
// the end user.
func BuildFailedResponse(message string, err any, data any) Response {
    return Response {
        Message: message,
        Error: err,
        Data: data,
    }
}
