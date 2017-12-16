# Versus
Versus is a project for comparing the performance of various languages.

## Perf Results
### 1000 Threads - Simple JSON Response
[Result](/tests/jmeter-results/simple-json-response/result.md)

## Setup Project Mac

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

## Running The Tests
Open Jmeter, then File->Open one of the jmx tests in the tests dir.
```
jmeter
```

