package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type (
	Config struct {
		Port string
	}
	SimpleResponse struct{
		Hello string `json:"hello""`
	}
)

func main() {
	fmt.Println("go app is running")
	startServer(getConfigFromEnvVariables())
}

func startServer(config Config) {
	fmt.Println("starting server with config: ", config)
	port := ":" + config.Port
	http.HandleFunc("/simple-json-response", simpleJsonResponse) // set router
	err := http.ListenAndServe(port, nil)                        // set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func simpleJsonResponse(response http.ResponseWriter, request *http.Request) {
	flusher, _ := response.(http.Flusher)
	response.Header().Add("Content-Type", "application/json")
	response.Header().Add("Connection", "keep-alive")  //node does this by default
	simpleResponse := SimpleResponse{Hello:"world"}
	json.NewEncoder(response).Encode(simpleResponse)
	flusher.Flush() //transfer encoding chunked. node does this by default
}

func getConfigFromEnvVariables() Config {
	return Config{
		Port: os.Getenv("PORT"),
	}
}
