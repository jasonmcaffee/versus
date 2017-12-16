const http = require('http');
const url = require('url');
const cluster = require('cluster');
const os = require('os');

function main(){
  console.log(`node app is running`);
  startServer();
}

function startServer({config=getConfigFromEnvVariables(), cpus=os.cpus()}={}){
  console.log(`starting server with config: ${JSON.stringify(config)}`);
  const {port, useCluster} = config;
  const shouldUseCluster = useCluster === 'true';
  if(shouldUseCluster){
    if (cluster.isMaster) {
      cpus.forEach(()=>cluster.fork());
    } else {
      createServerAndListen({port});
    }
  }else{
    createServerAndListen({port});
  }
}
function createServerAndListen({port}){
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
    port: process.env.PORT,
    useCluster: process.env.USE_CLUSTER,
  }
}

main();