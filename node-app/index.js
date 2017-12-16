const http = require('http');

function main(){
  console.log(`node app is running`);
  startServer();
}

function startServer(){
  const server = http.createServer((req, res) => {
    res.writeHead(200, {'Content-type':'text/plan'});
    res.write('Hello Node JS Server Response');
    res.end();
  });

  server.listen(8000);
}

main();