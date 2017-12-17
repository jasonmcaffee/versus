package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"io/ioutil"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
)

type (
	Config struct {
		Port string
		DbHost string
		DbUser string
		DbPassword string
		DbPort string
		DbConnectionLimit int
		DbSchema string
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
	DbOperationsRequest struct {
		IntColumn int `json:"intColumn"`
		StringColumn string `json:"stringColumn"`
	}
)

func main() {
	fmt.Println("go app is running")
	startServer(getConfigFromEnvVariables())
}


//########################################################################################################## test 1
func simpleJsonResponse(response http.ResponseWriter, request *http.Request) {
	flusher, _ := response.(http.Flusher)
	response.Header().Add("Content-Type", "application/json")
	response.Header().Add("Connection", "keep-alive")  //node does this by default
	simpleResponse := SimpleResponse{Hello:"world"}
	json.NewEncoder(response).Encode(simpleResponse)
	flusher.Flush() //transfer encoding chunked. node does this by default
}

//########################################################################################################## test 2
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

//########################################################################################################## test 3
func dbOperations(response http.ResponseWriter, request *http.Request) {
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
	jsonObject := &DbOperationsRequest{}
	err = json.Unmarshal(b, jsonObject)
	if err != nil {
		http.Error(response, err.Error(), 500)
		return
	}

	conn := getDbConnection()
	query := "select 1+1"
	rows := dbQuery(conn, query)

	sendJsonResponse(response, rows)
}

func dbQuery(conn *sql.DB, query string) *sql.Rows{
	results, err := conn.Query(query)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	return results
}

var dbConnection *sql.DB
func getDbConnection() *sql.DB{
	if dbConnection != nil{
		return dbConnection
	}
	config := getConfigFromEnvVariables()
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", config.DbUser, config.DbPassword, config.DbHost, config.DbPort, config.DbSchema)
	dbConnection, err := sql.Open("mysql", connectionString)
	dbConnection.SetMaxIdleConns(config.DbConnectionLimit)//go closes connections quickly, so force them to stay open.
	dbConnection.SetMaxOpenConns(config.DbConnectionLimit)
	if err != nil {
		panic(err.Error())
	}
	return dbConnection
}

//########################################################################################################## common
func startServer(config Config) {
	fmt.Println("starting server with config: ", config)
	port := ":" + config.Port
	http.HandleFunc("/simple-json-response", simpleJsonResponse) // set router
	http.HandleFunc("/accept-and-return-json", acceptAndReturnJson)
	http.HandleFunc("/db-operations", dbOperations)
	err := http.ListenAndServe(port, nil)                        // set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func sendJsonResponse(response http.ResponseWriter, data interface{}) {
	flusher, _ := response.(http.Flusher)
	response.Header().Add("Content-Type", "application/json")
	response.Header().Add("Connection", "keep-alive")  //node does this by default
	json.NewEncoder(response).Encode(data)
	flusher.Flush() //transfer encoding chunked. node does this by default
}

func getConfigFromEnvVariables() Config {
	dbConnectionLimit, _ := strconv.Atoi( os.Getenv("DB_CONNECTION_LIMIT") )
	return Config{
		Port: os.Getenv("PORT"),
		DbHost: os.Getenv("DB_HOST"),
		DbUser: os.Getenv("DB_USER"),
		DbPassword: os.Getenv("DB_PASSWORD"),
		DbPort: os.Getenv("DB_PORT"),
		DbSchema: os.Getenv("DB_SCHEMA"),
		DbConnectionLimit: dbConnectionLimit,
	}
}
