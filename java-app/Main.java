import java.io.IOException;
import java.io.OutputStream;
import java.net.InetSocketAddress;

import com.sun.net.httpserver.HttpExchange;
import com.sun.net.httpserver.HttpHandler;
import com.sun.net.httpserver.HttpServer;
import com.sun.net.httpserver.Headers;

public class Main {

    public static void main(String[] args) throws Exception {
        startServer(new Config());
    }

    public static void startServer(Config config) throws Exception {
        System.out.println("starting server with config: " + config);
        int port = config.getPort();
        HttpServer server = HttpServer.create(new InetSocketAddress(port), 0);
        server.createContext("/simple-json-response", new SimpleJsonResponseHandler());
        server.setExecutor(null); // creates a default executor
        server.start();
    }

    static class SimpleJsonResponseHandler implements HttpHandler {
        @Override
        public void handle(HttpExchange t) throws IOException {
            SimpleJsonResponse simpleJsonResponse = new SimpleJsonResponse();
            simpleJsonResponse.setHello("world");
            String response = simpleJsonResponse.toJson();

            Headers headers = t.getResponseHeaders();
            headers.add("Content-Type", "application/json");
            headers.add("Connection", "keep-alive");

            t.sendResponseHeaders(200, response.length());

            OutputStream os = t.getResponseBody();
            os.write(response.getBytes());
            os.close();
        }
    }

    static class SimpleJsonResponse {
        private String hello;
        public void setHello(String val){
            this.hello = val;
        }
        public String toJson(){
            String s = "{ \"hello\": \"" + this.hello + "\"}";
            return s;
        }
    }

    static class Config {
        public int getPort(){
            String portString = System.getenv("PORT");
            int port = Integer.parseInt(portString);
            return port;
        }

        @Override
        public String toString(){
            String s = "Port: " + this.getPort();
            return s;
        }
    }
}