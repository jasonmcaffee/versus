const http = require('http');
const url = require('url');
const cluster = require('cluster');
const os = require('os');
const mysql = require('mysql');

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
  let query =  `select 1 + 1`;
  let rows = await dbQuery({query});
  sendJsonResponse(rows, response);
}

function dbQuery({conn=getDbConnection(), query}){
  return new Promise((resolve, reject)=>{
    conn.query(query, function(err, results, fields) {
      if(err){return reject(err);}
      resolve(results);
    });
  });
}

let dbConnection;
function getDbConnection({config=getConfigFromEnvVariables()}={}){
  dbConnection = dbConnection || mysql.createPool({
    connectionLimit : config.dbConnectionLimit,
    host     : config.dbHost,
    user     : config.dbUser,
    password : config.dbPassword,
    database : config.dbSchema,
    port: config.dbPort,
    debug    :  false
  });
  return dbConnection;
}


//########################################################################################################## common
function sendJsonResponse(json, response){
  const headers = {'Content-Type': 'application/json'};
  const simpleResponse = {hello: 'world'};
  response.writeHead(200, headers);
  response.end(JSON.stringify(json));
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
    default:
      notFoundReponse(request, response);
  }
}

main();