# Services app

Boiler plate golang micro services app built with Gin, ElasticSearch, MySQL, Cobra, Viper, Docker-compose

### Features
- [X] Modular request flow - 
  - [X] handlers
  - [X] validators
  - [X] usecases
  - [X] entity
  - [ ] serialiser
  - [X] response
- [X] Authentication
	- [ ] Social auth
	- [X] lock on exceeding maximum attempts (5) - configurable
- [ ] Notifications - (NSQ)
	- [ ] Email
	- [ ] SMS
- [ ] Cache API
	- [ ] Redis
- [X] Logging API - Logrus
- [X] Unit + Integration testcase
- [X] Code coverage
- [ ] CI - Travis CI

### Tools
- [X] Cobra command cli apps
- [X] Viper configuration management
- [X] ElasticSearch 
	- [X] CRUD API
	- [X] flexible search querying using [elastic]
	- [ ] Aggregate queries
	- [ ] Geo distance queries
	- [X] Index management using JSON
- [ ] MySQL
	- [ ] CRUD API
	- [ ] Install scripts
	- [ ] Management scripts - backup, Point in time recovery
	- [ ] Monitoring API
- [X] Makefile commands 
	- [X] build
	- [X] install
	- [X] run
	- [X] format
	- [X] lint
	- [X] vet
	- [X] test
- [ ] Postman collections
- [ ] Distributed task queue (NSQ)
- [ ] Logstash + filebeats setup for log filtering and analysis

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
make setup
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


[//]: # (reference links here)

[elastic]: <https://olivere.github.io/elastic/>