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
	DbOperationsResult struct {
		ID int `json:"id"`
		IntColumn int `json:"intColumn"`
		StringColumn string `json:"stringColumn"`
	}
)

func main() {
	//runtime.GOMAXPROCS(runtime.NumCPU())
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
	//insert
	insertQuery :="insert into db_operations (stringColumn, intColumn) values (?, ?)"
	_, lastInsertId := dbUpdate(conn, insertQuery, jsonObject.StringColumn, jsonObject.IntColumn)

	//read
	query := "select * from db_operations where id = ?"
	rows := dbQuery(conn, query, lastInsertId)

	//delete
	deleteQuery := "delete from db_operations where id = ?"
	_, _ = dbUpdate(conn, deleteQuery, lastInsertId)

	//return result
	result := []DbOperationsResult{}
	for rows.Next() {
		var dbOperationsResult DbOperationsResult
		err = rows.Scan(&dbOperationsResult.ID, &dbOperationsResult.StringColumn, &dbOperationsResult.IntColumn)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		result = append(result, dbOperationsResult)
	}

	rows.Close()
	sendJsonResponse(response, result)
}

func dbQuery(conn *sql.DB, query string, args ...interface{}) *sql.Rows{
	results, err := conn.Query(query, args...)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	return results
}

func dbUpdate(conn *sql.DB, query string, args ...interface{}) (result sql.Result, lastInsertId int64){
	result, err := conn.Exec(query, args...)
	if err != nil{
		panic(err.Error())
	}
	lastInsertId, err = result.LastInsertId()
	if err != nil{
		panic(err.Error())
	}
	return result, lastInsertId
}

var dbConnection *sql.DB
func getDbConnection() *sql.DB{
	if dbConnection != nil{
		return dbConnection
	}
	config := getConfigFromEnvVariables()
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", config.DbUser, config.DbPassword, config.DbHost, config.DbPort, config.DbSchema)
	conn, err := sql.Open("mysql", connectionString)
	dbConnection = conn
	dbConnection.SetMaxIdleConns(config.DbConnectionLimit)
	//dbConnection.SetConnMaxLifetime(time.Second * 1)
	dbConnection.SetMaxOpenConns(config.DbConnectionLimit)
	if err != nil {
		if dbConnection != nil{
			dbConnection.Close()
		}
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
