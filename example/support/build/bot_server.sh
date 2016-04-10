#!/bin/sh

cd `dirname $0`
#echo "start to embed resources"
#cd ../../resource
#rice embed-go

echo "start to build 'bot_server'"
cd ../../app/bot_server
GOOS=linux GOARCH=amd64 go build 
mv bot_server ../../support/deploy/roles/deploy-app/files/
echo "moved to support/deploy/roles/deploy-app/files/bot_server"
