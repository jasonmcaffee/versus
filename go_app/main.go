package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"io/ioutil"
)

type (
	Config struct {
		Port string
	}
	SimpleResponse struct{
		Hello string `json:"hello"`
	}
	AcceptAndReturnJsonRequest struct{
		String string `json:"string"`
		Number int `json:"number"`
		Boolean bool `json:"boolean"`
		ArrayNumber []int `json:"array number"`
		ArrayString []string `json:"array string"`
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
	http.HandleFunc("/accept-and-return-json", acceptAndReturnJson)
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

//https://blog.golang.org/json-and-go
func acceptAndReturnJson(response http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		http.Error(response, "not found", 404)
		return
	}
	b, err := ioutil.ReadAll(request.Body)
	defer request.Body.Close()
	if err != nil {
		http.Error(response, err.Error(), 500)
		return
	}
	jsonObject := &AcceptAndReturnJsonRequest{}
	err = json.Unmarshal(b, jsonObject)
	if err != nil {
		http.Error(response, err.Error(), 500)
		return
	}
	sendJsonResponse(response, *jsonObject)
}

func sendJsonResponse(response http.ResponseWriter, data interface{}) {
	flusher, _ := response.(http.Flusher)
	response.Header().Add("Content-Type", "application/json")
	response.Header().Add("Connection", "keep-alive")  //node does this by default
	json.NewEncoder(response).Encode(data)
	flusher.Flush() //transfer encoding chunked. node does this by default
}

func getConfigFromEnvVariables() Config {
	return Config{
		Port: os.Getenv("PORT"),
	}
}
