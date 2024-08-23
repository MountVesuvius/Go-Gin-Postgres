package helpers

type Response struct {
	Message string `json:"message"`
	Error   any    `json:"error,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func BuildSuccessfulResponse(message string, data any) Response {
    return Response {
        Message: message,
        Data: data,
    }
}

func BuildFailedResponse(message string, err any, data any) Response {
    return Response {
        Message: message,
        Error: err,
        Data: data,
    }
}
