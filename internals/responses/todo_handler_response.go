package responses

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func SuccessJSONResponse(writer http.ResponseWriter, code int, data interface{}) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(code)
	if data != nil {
		if err := json.NewEncoder(writer).Encode(data); err != nil {
			http.Error(writer, "Error encoding response", http.StatusInternalServerError)
		}
	}
}

func ErrorJSONResponse(writer http.ResponseWriter, code int, message string) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(code)
	if err := json.NewEncoder(writer).Encode(map[string]string{"error": message}); err != nil {
		fmt.Println(err)
		http.Error(writer, "Error encoding response", http.StatusInternalServerError)
	}
}
