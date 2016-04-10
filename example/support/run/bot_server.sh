#!/bin/sh

cd `dirname $0`
#cd ../../resource
##rice embed-go
cd ../../app/bot_server
go run main.go -debug=true -config=../../support/config/develop.toml
