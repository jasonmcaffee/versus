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
func performHttpRequest(response http.ResponseWriter, request *http.Request) {
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

	path := "http://localhost:7878/accept-and-return-json"
	method := "POST"
	headers := createCommonHeaders()
	responseObject := &AcceptAndReturnJsonRequest{}
	_, err = req(path, method, headers, jsonObject, responseObject)

	sendJsonResponse(response, responseObject)
}

```
![Summary](go-summary.png)

![Response](go-response-times.png)

## Node with Cluster
```js
async function performHttpRequest(request, response){
  let data = await getJsonRequest(request);
  let result = await req({hostname:'localhost', port:7878, path:'/accept-and-return-json', method:'POST', data});
  sendJsonResponse(result, response);
}
```
![Summary](nodecluster-summary.png)

![Response](nodecluster-response-times.png)

## Node
```js
async function performHttpRequest(request, response){
  let data = await getJsonRequest(request);
  let result = await req({hostname:'localhost', port:7878, path:'/accept-and-return-json', method:'POST', data});
  sendJsonResponse(result, response);
}
```
![Summary](node-summary.png)

![Response](node-response-times.png)


