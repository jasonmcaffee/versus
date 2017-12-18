const http = require('http');
const url = require('url');
const cluster = require('cluster');
const os = require('os');
const mysql = require('mysql2');

function main(){
  console.log(`node app is running`);
  startServer();
}

//########################################################################################################## test 1
function simpleJsonResponse(request, response){
  const headers = {'Content-Type': 'application/json'};
  const simpleResponse = {hello: 'world'};
  response.writeHead(200, headers);
  response.end(JSON.stringify(simpleResponse));
}

//########################################################################################################## test 2
function acceptAndReturnJson(request, response){
  if (request.method != 'POST'){ return notFoundReponse(request, response); }
  let body = '';
  request.on('data', function(data){
    body += data;
  });
  request.on('end', function(data){
    let json = JSON.parse(body);
    sendJsonResponse(json, response);
  });
}

//########################################################################################################## test 3
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

function dbQuery({conn=getDbConnection(), query, data}){
  return new Promise((resolve, reject)=>{
    conn.query(query, data, function(err, results, fields) {
      if(err){return reject(err);}
      resolve(results);
    });
  });
}

let dbConnection;
function getDbConnection({config=getConfigFromEnvVariables(), cpus=os.cpus()}={}){
  //since each process in the cluster will start a process, we need to ensure we are using the desired max connection limit
  let dbConnectionLimit = config.useCluster === 'true' ? config.dbConnectionLimit / cpus.length : config.dbConnectionLimit;
  dbConnection = dbConnection || mysql.createPool({
    connectionLimit : dbConnectionLimit,
    host     : config.dbHost,
    user     : config.dbUser,
    password : config.dbPassword,
    database : config.dbSchema,
    port: config.dbPort,
    debug    :  false
  });
  return dbConnection;
}

//########################################################################################################## test 4
async function performHttpRequest(request, response){
  let data = await getJsonRequest(request);
  let result = await req({hostname:'localhost', port:7878, path:'/accept-and-return-json', method:'POST', data});
  sendJsonResponse(result, response);
}

function req({method='GET', data='', path='/path', port=80, contentType='application/json', hostname='www.google.com', agent=getHttpAgent()}){
  data = JSON.stringify(data);
  return new Promise((resolve, reject)=>{
    const options = { hostname, port, path, method, agent,
      headers: {
        'Content-Type': contentType,
        'Content-Length': Buffer.byteLength(data)
      }
    };
    let req = http.request(options, (res)=>{
      let body = '';
      res.on('data', (chunk) => {
        body += chunk;
      });
      res.on('end', () => {
        let result = JSON.parse(body);
        resolve(result);
      });
    });
    req.on('error', (e)=>{
      reject(e);
    });
    if(data != ''){
      req.write(data);
    }
    req.end();
  });
}

let httpAgent;
function getHttpAgent(){
  if(httpAgent){
    return httpAgent
  }
  let config = getConfigFromEnvVariables();
  httpAgent = new http.Agent({
    keepAlive: true,
    maxSockets: config.httpRequestSockets,
    maxFreeSockets: config.httpRequestSockets,
  });
  return httpAgent;
}

//########################################################################################################## test 5
async function findPrimeNumbers(request, response){
  let requestObject = await getJsonRequest(request);
  let min = requestObject.min;
  let max = requestObject.max;
  let primesArray = getPrimeNumbersBetween(min, max);
  let result = {numberOfPrimes:primesArray.length};
  sendJsonResponse(result, response);
}

function isPrime(num) {
  const sqrtnum=Math.floor(Math.sqrt(num)) + 1;
  let prime = num != 1;
  for(let i = 2; i < sqrtnum; i++) {
    if(num % i == 0) {
      prime = false;
      break;
    }
  }
  return prime;
}

function getPrimeNumbersBetween(min, max){
  let primes = [];
  for(let i = min; i <= max; ++i){
    if(isPrime(i)){
      primes.push(i);
    }
  }
  return primes;
}


//########################################################################################################## common
function sendJsonResponse(json, response){
  const headers = {'Content-Type': 'application/json'};
  const simpleResponse = {hello: 'world'};
  response.writeHead(200, headers);
  response.end(JSON.stringify(json));
}

function getJsonRequest(request){
  return new Promise((resolve, reject)=>{
    let body = '';
    request.on('data', function(data){
      body += data;
    });
    request.on('end', function(data){
      let json = JSON.parse(body);
      resolve(json);
    });
  });
}

function notFoundReponse(request, response){
  response.writeHead(404);
  response.end();
}

function getConfigFromEnvVariables(){
  return {
    port: process.env.PORT,
    useCluster: process.env.USE_CLUSTER,
    dbUser: process.env.DB_USER,
    dbPassword: process.env.DB_PASSWORD,
    dbHost: process.env.DB_HOST,
    dbSchema: process.env.DB_SCHEMA,
    dbConnectionLimit: parseInt(process.env.DB_CONNECTION_LIMIT),
    dbPort: parseInt(process.env.DB_PORT),
    httpRequestSockets: parseInt(process.env.HTTP_REQUEST_SOCKETS),
  };
}

function startServer({config=getConfigFromEnvVariables(), cpus=os.cpus()}={}){
  console.log(`starting server with config: ${JSON.stringify(config)}`);
  const {port, useCluster} = config;
  const shouldUseCluster = useCluster === 'true';
  if(shouldUseCluster){
    if (cluster.isMaster){
      cpus.forEach(()=>cluster.fork());
    }else{
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
    case '/accept-and-return-json':
      acceptAndReturnJson(request, response);
      break;
    case '/db-operations':
      dbOperations(request, response);
      break;
    case '/perform-http-request':
      performHttpRequest(request, response);
      break;
    case '/find-prime-numbers':
      findPrimeNumbers(request, response);
      break;
    default:
      notFoundReponse(request, response);
  }
}

main();