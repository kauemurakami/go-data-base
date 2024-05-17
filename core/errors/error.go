package http_error

import "encoding/json"

func ResponseError(message string) []byte {
	errorMessage := map[string]string{"error": message}
	jsonerror, err := json.Marshal(errorMessage)
	if err != nil {
		errorMessage["error"] = "Erro ao converter dados"
		jsonerror, _ = json.Marshal(errorMessage)
	}
	return jsonerror
}
