#!/bin/sh
# https://github.com/azer/radio-paradise/blob/master/Makefile
# GOPATH=$(shell pwd)/vendor:$(shell pwd)
GOBIN=$(shell pwd)
GOFILES=$(wildcard *.go)
GONAME=$(shell basename "$(PWD)")
LINT_FILES := $(shell find . -name '*.go' | grep -v /vendor/)
PROJECT_ROOT="$(PWD)"
PID_FILE=$(GONAME).pid
PID=`cat $(PROJECT_ROOT)/$(PID_FILE)`
port=3333

build: fmt lint vet
	@echo "Building $(GOFILES) to ./"
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go build -o $(GONAME) $(GOFILES)

get:
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go get . ; dep ensure;

install:
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go install $(GOFILES)

fmt:
	@for file in ${LINT_FILES} ;  do \
		go fmt $$file ; \
	done

lint:
	@for file in ${LINT_FILES} ;  do \
		golint $$file ; \
	done

vet:
	@for file in ${LINT_FILES} ;  do \
		go tool vet $$file ; \
	done	

run:
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go run $(GOFILES) serve -p=$(port)

watch: build stop start
	@fswatch -o *.go **/*.go | xargs -n1 -I{}  make restart || make stop

restart: clear stop clean build start

start: clear build
	@echo "Starting ./$(GONAME)"
	@./$(GONAME) serve -p=$(port) & echo $$! > $(PID_FILE)

stop:
	@echo "Stopping ./$(GONAME) if it's running"
	@if [ -f $(PROJECT_ROOT)/$(PID_FILE) ] ; \
	then \
		echo "Killing services process: "$(PID) ; \
	    kill -9 $(PID) ; \
	fi;

clear:
	@clear

clean:
	@echo "Cleaning"
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go clean

.PHONY: build get install run watch start stop restart clean