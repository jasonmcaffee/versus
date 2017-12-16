h1. Versus
Versus is a project for comparing the performance of various languages.

h2. Setup Project Mac

h3. Environment Variables
h4. PORT
port for server to listen on
h4. USE_CLUSTER
indicates whether node should use a cluster to utilize all cpus

h2. Setup Software Mac
h3. JMeter
Jmeter requires java version 8 (9 doesn't work)
```
brew cask install caskroom/versions/java8
brew install jmeter --with-plugins
```

h3. Node
install nvm and node 8.5

```
touch ~/.profile
curl -o- https://raw.githubusercontent.com/creationix/nvm/v0.32.1/install.sh | bash
. ~/.nvm/nvm.sh
source ~/.profile
nvm install 8.5.0
nvm use 8.5.0
```

h2. Running The Tests
Open Jmeter
```
jmeter
```

