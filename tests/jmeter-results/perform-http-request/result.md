# Perf Results
This performance tests sees how fast apps can:
- accept json post
- perform http post with json to accept-and-return-json endpoint
- read result for accept-and-return-json
- send json response with result

Each app was started, then given a few test runs before the results were recorded.

JMeter was used with 1000 threads for the thread group.

Each app was configured to use a http connection pool with a max size of 100.

Node was tested with clustering turned on and off.

Performance between Node and Go appears to be relatively the same.

## Go
```go
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

```
![Summary](go-summary.png)

![Response](go-response-times.png)

## Node with Cluster
```js
async function dbOperations(request, response){
  let jsonRequest = await getJsonRequest(request);
  //insert
  let insertQuery = `insert into db_operations set ?`;
  let insertRows = await dbQuery({query:insertQuery, data:jsonRequest});
  let insertId = insertRows.insertId;
  //retrieve the inserted row
  let query = `select * from db_operations where id = ${insertId}`;
  let result = await dbQuery({query});
  //delete the inserted row
  let deleteQuery = `delete from db_operations where id = ${insertId}`;
  await dbQuery({query:deleteQuery});
  sendJsonResponse(result, response);
}
```
![Summary](nodecluster-summary.png)

![Response](nodecluster-response-times.png)

## Node
```js
async function dbOperations(request, response){
  let jsonRequest = await getJsonRequest(request);
  //insert
  let insertQuery = `insert into db_operations set ?`;
  let insertRows = await dbQuery({query:insertQuery, data:jsonRequest});
  let insertId = insertRows.insertId;
  //retrieve the inserted row
  let query = `select * from db_operations where id = ${insertId}`;
  let result = await dbQuery({query});
  //delete the inserted row
  let deleteQuery = `delete from db_operations where id = ${insertId}`;
  await dbQuery({query:deleteQuery});
  sendJsonResponse(result, response);
}

```
![Summary](node-summary.png)

![Response](node-response-times.png)


