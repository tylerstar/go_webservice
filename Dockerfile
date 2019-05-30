FROM golang:latest

MAINTAINER Tyler Leung binglux@gmail.com

WORKDIR /app
ADD bin/webservice_linux_amd64 /app

