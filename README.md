# Versus
Versus is a project for comparing the performance of various languages.

## Perf Results

### 1000 Threads - Simple JSON Response
Vanilla servers with endpoint to return a json response from a GET request.

Node and Go perform relatively the same at requests handled per second.

[Result](/tests/jmeter-results/simple-json-response/result.md)

### 1000 Threads - Simple JSON Response
Vanilla servers with endpoint to accept and return a json object from a POST request.

Node and Go perform relatively the same at requests handled per second.
[Result](/tests/jmeter-results/accept-and-return-json/result.md)

## Setup Project Mac
### ulimit
By default, the max amount of file descriptors is set to a low number, and this will affect the max number of connections the apps can handle.

Fix this by increasing the ulimit.

Check current settings
```
sudo launchctl limit
```

Update settings to unlimited.
```
echo "limit maxfiles 1024 unlimited" | sudo tee -a /etc/launchd.conf
```

[Stack Overflow For Help With Ulimit](https://superuser.com/questions/302754/increase-the-maximum-number-of-open-file-descriptors-in-snow-leopard)
### Environment Variables
#### PORT
port for server to listen on
#### USE_CLUSTER
indicates whether node should use a cluster to utilize all cpus

## Setup Software Mac
### JMeter
Jmeter requires java version 8 (9 doesn't work)
```
brew cask install caskroom/versions/java8
brew install jmeter --with-plugins
```
You may need to increase the heap size if you see jmeter throw memory exceptions.
Open
```
/usr/local/Cellar/jmeter/3.3/bin/jmeter
```
And add the following to the top of the file
```
#!/bin/bash
HEAP="-Xms1024m -Xmx2048m"
```
### Node
install nvm and node 8.5

```
touch ~/.profile
curl -o- https://raw.githubusercontent.com/creationix/nvm/v0.32.1/install.sh | bash
. ~/.nvm/nvm.sh
source ~/.profile
nvm install 8.5.0
nvm use 8.5.0
```

### Go
install go
```
brew install go
```

add the following to your ~/.bash_profile
```
export GOPATH=/Users/jason.mcaffee/Documents/dev/go
export GOROOT=/usr/local/Cellar/go/1.8/libexec
export PATH=$PATH:$GOPATH/bin
```

make sure this project lives in a directory inside the GOPATH directory

## Running The Tests
Open Jmeter, then File->Open one of the jmx tests in the tests dir.

To run jmeter, open a terminal and run
```
jmeter
```

A UI should pop up, where you can run your tests.

