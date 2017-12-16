# Perf Results
This performance test sees how many requests standard go and node apps (no external frameworks or libs) can serve a simple json response from a GET endpoint.

Each app was started, then given a few test runs before the results were recorded.

JMeter was used with 1000 threads for the thread group.

Node was tested with clustering turned on and off.

Performance appears to be relatively the same.

## Go
```go
func simpleJsonResponse(response http.ResponseWriter, request *http.Request) {
	flusher, _ := response.(http.Flusher)
	response.Header().Add("Content-Type", "application/json")
	response.Header().Add("Connection", "keep-alive")  //node does this by default
	simpleResponse := SimpleResponse{Hello:"world"}
	json.NewEncoder(response).Encode(simpleResponse)
	flusher.Flush() //transfer encoding chunked. node does this by default
}
```
![Summary](go-summary.png)

![Response](go-response-times.png)

## Node with Cluster
```js
 function simpleJsonResponse(request, response){
   const headers = {'Content-Type': 'application/json'};
   const simpleResponse = {hello: 'world'};
   response.writeHead(200, headers);
   response.end(JSON.stringify(simpleResponse));
 }
```
![Summary](nodecluster-summary.png)

![Response](nodecluster-response-times.png)

## Node
```js
 function simpleJsonResponse(request, response){
   const headers = {'Content-Type': 'application/json'};
   const simpleResponse = {hello: 'world'};
   response.writeHead(200, headers);
   response.end(JSON.stringify(simpleResponse));
 }
```
![Summary](node-summary.png)

![Response](node-response-times.png)

