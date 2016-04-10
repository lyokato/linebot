#!/bin/sh

cd `dirname $0`
ansible-playbook bot_server.yml -i hosts -u centos --private-key=$EC2_KEY
