#!/bin/bash

while true
do

ts=`date +%s`;
curl -X POST -d "[{\"metric\": \"test-metric\", \"endpoint\": \"test-endpoint1\", \"timestamp\": $ts,\"step\": 60,\"value\": 1,\"counterType\": \"GAUGE\",\"tags\": \"idc=lg,project=xx\"}]" http://127.0.0.1:1988/v1/push

ts=`date +%s`;
curl -X POST -d "[{\"metric\": \"test-metric\", \"endpoint\": \"test-endpoint2\", \"timestamp\": $ts,\"step\": 60,\"value\": 2,\"counterType\": \"GAUGE\",\"tags\": \"idc=lg,project=xx\"}]" http://127.0.0.1:1988/v1/push

ts=`date +%s`;
curl -X POST -d "[{\"metric\": \"test-metric\", \"endpoint\": \"test-endpoint3\", \"timestamp\": $ts,\"step\": 60,\"value\": 3,\"counterType\": \"GAUGE\",\"tags\": \"idc=lg,project=xx\"}]" http://127.0.0.1:1988/v1/push
sleep 60
done
