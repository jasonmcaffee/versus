const http = require('http');

function main(){
  console.log(`node app is running`);
  startServer();
}

function startServer({config=getConfigFromEnvVariables()}={}){
  console.log(`starting server with config: ${JSON.stringify(config)}`);

  const {port} = config;
  const server = http.createServer((req, res) => {
    res.writeHead(200, {'Content-type':'text/plan'});
    res.write('Hello Node JS Server Response');
    res.end();
  });

  server.listen(port);
}

function getConfigFromEnvVariables(){
  return {
    port: process.env.PORT
  }
}

main();