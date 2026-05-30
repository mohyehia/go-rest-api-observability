package health

import (
	"log"
	"net/http"
)

func Handler(writer http.ResponseWriter, _ *http.Request) {
	log.Println("health check :: start")
	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	_, err := writer.Write([]byte(`{"status": "UP"}`))
	if err != nil {
		log.Printf("Error writing health check response: %v\n", err)
	}
	log.Println("health check :: end")
}
