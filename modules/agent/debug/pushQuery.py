#!/usr/bin/env python
#!-*- coding:utf8 -*-

import requests
import time
import json

ts = int(time.time())
payload = [
    {
        "endpoint": "test-endpoint",
        "metric": "mike",
        "timestamp": ts,
        "step": 60,
        "value": 1,
        "counterType": "GAUGE",
        "tags": "tag1=lg,tag2=beijing",
        "fields": "field1=This metric will not enter influxdb because conterType"
    },
    {
        "endpoint": "test-endpoint",
        "metric": "mike",
        "timestamp": ts,
        "step": 60,
        "value": 12345,
        "counterType": "MQ,KEYWORD",
        "tags": "tag3=lg,tag4=beijing",
        "fields": "field1=Mike does not have a 170mm penis!, field2=7777"
    },
    {
        "endpoint": "test-endpoint",
        "metric": "mike",
        "timestamp": ts,
        "step": 60,
        "value": 678910,
        "counterType": "MQ,KEYWORD",
        "tags": "tag3=lg,tag4=beijing",
        "fields": "field1=Mike does not have a 175mm penis!, field2=This metric will overwrite previous value because the (metric - tags - endpoint) are the same."
    },
    {
        "endpoint": "test-endpoint-2",
        "metric": "mike",
        "timestamp": ts,
        "step": 60,
        "value": 12345,
        "counterType": "MQ,KEYWORD",
        "tags": "tag3=lg,tag4=beijing",
        "fields": "field1=Mike does not have a 170mm penis!, field2=This metric will not overwrite."
    },
]

r = requests.post("http://10.20.30.40:1988/v1/push", data=json.dumps(payload))

print r.text
