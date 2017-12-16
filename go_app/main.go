package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"os"
)

type(
	Config struct{
		Port string
	}
)

func main() {
	fmt.Println("go app is running")
	startServer(getConfigFromEnvVariables())
}

func startServer(config Config) {
	fmt.Println("starting server with config: ", config)
	port := ":" +config.Port
	http.HandleFunc("/", sayhelloName)       // set router
	err := http.ListenAndServe(port, nil) // set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func getConfigFromEnvVariables() Config {
	return Config{
		Port: os.Getenv("PORT"),
	}
}

func sayhelloName(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()       // parse arguments, you have to call this by yourself
	fmt.Println(r.Form) // print form information in server side
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	fmt.Fprintf(w, "Hello astaxie!") // send data to client side
}
