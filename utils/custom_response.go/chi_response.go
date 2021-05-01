package custom_response

import (
	"encoding/json"
	"log"
	"net/http"
)

type CustomResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func ErrorResponse(w http.ResponseWriter, statusCode int, errorMsg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := CustomResponse{
		Message: "Error",
		Data:    errorMsg,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Fatal(err)
	}
}

func SuccessResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := CustomResponse{
		Message: "Success",
		Data:    data,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Fatal(err)
	}
}
