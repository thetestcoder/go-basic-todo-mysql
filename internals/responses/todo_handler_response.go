package responses

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type StandardResponse struct {
	Success    bool        `json:"success"`
	Message    string      `json:"message,omitempty"`
	Data       interface{} `json:"data,omitempty"`
	Error      string      `json:"error,omitempty"`
	StackTrace interface{} `json:"stacktrace,omitempty"`
}

func SuccessJSONResponse(writer http.ResponseWriter, code int, data interface{}) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(code)

	response := StandardResponse{
		Success: true,
		Data:    data,
	}

	if data != nil {
		if err := json.NewEncoder(writer).Encode(response); err != nil {
			http.Error(writer, "Error encoding response", http.StatusInternalServerError)
		}
	}
}

func ErrorJSONResponse(writer http.ResponseWriter, code int, message string, stackTrace interface{}) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(code)

	response := StandardResponse{
		Success:    false,
		Error:      message,
		StackTrace: stackTrace,
	}

	if err := json.NewEncoder(writer).Encode(response); err != nil {
		fmt.Println(err)
		http.Error(writer, "Error encoding response", http.StatusInternalServerError)
	}
}
