# Services app

Boiler plate golang micro services app built with Gin, ElasticSearch, MySQL, Cobra, Viper, Docker-compose

## Pre-requisites

1. MySQL should be installed 
2. Docker-compose should be installed
3. Golang should be installed
4. Golang dependency manager : Dep
5. Golang lint tool
> golang lint dependency for makefile building
> go get -u golang.org/x/lint/golint

Once the above requirements are installed, we should be able to get 
the app running in a breeze.

> Note: Most of the commands below need to be run from GIT repo directory.

### Steps to setup
1. Increase memory for ElasticSearch docker

```
sudo sysctl -w vm.max_map_count=262144
```
2. Start ElasticSearch and Kibana using docker-compose

```
# starts elastic search (9200) and kibana (5601)
docker-compose -f docker-compose-dev.yml up 
```
3. The app uses Makefile heavily for build, test and deploy. Makefile to get the dependencies for the app
```
# get golang dependencies
make get
```
4. Run the app
```
make run
```

### Code formats
Go lint tool
```
make lint
```
Go format tool
```
make fmt
```
Go vet tool
```
make vet
```
Go build
```
make build
```
Go install
```
make install
```
