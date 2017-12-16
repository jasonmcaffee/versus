const http = require('http');
const url = require('url');

function main(){
  console.log(`node app is running`);
  startServer();
}

function startServer({config=getConfigFromEnvVariables()}={}){
  console.log(`starting server with config: ${JSON.stringify(config)}`);
  const {port} = config;
  const server = http.createServer((request, response) => {
    route(request, response);
  });
  server.listen(port);
}

function route(request, response){
  const urlParts = url.parse(request.url);
  const path = urlParts.pathname;
  switch(path){
    case '/simple-json-response':
      simpleJsonResponse(request, response);
      break;
    default:
      notFoundReponse(request, response);
  }
}

function simpleJsonResponse(request, response){
  const headers = {'Content-Type': 'application/json'};
  const simpleResponse = {hello: 'world'};
  response.writeHead(200, headers);
  response.end(JSON.stringify(simpleResponse));
}

function notFoundReponse(request, response){
  response.writeHead(404);
  response.end();
}

function getConfigFromEnvVariables(){
  return {
    port: process.env.PORT
  }
}

main();