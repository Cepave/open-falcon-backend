#!/usr/bin/env python
#!-*- coding:utf8 -*-

import requests
import time
import json

ts = int(time.time())
payload = [
    {
        "endpoint": "ctl-js-AAA",
        "metric": "mike",
        "timestamp": ts,
        "step": 60,
        "value": 6,
        "counterType": "GAUGE",
        "tags": "tag3=lg,tag4=beijing",
        "fields": "field1=This metric will not enter influxdb because conterType"
    },
    {
        "endpoint": "ctl-js-AAA",
        "metric": "service.http.lvs.443.port",
        "timestamp": ts,
        "step": 60,
        "value": 6,
        "counterType": "GAUGE",
        "tags": "",
        "fields": "field1=This metric will not enter influxdb because conterType"
    },
    {
        "endpoint": "ctl-js-AAA",
        "metric": "mike",
        "timestamp": ts,
        "step": 60,
        "value": 4,
        "counterType": "MQ,KEYWORD",
        "tags": "",
        "fields": "lvsvip=www.google.com, field2=This metric will be overwritten"
    },
    {
        "endpoint": "ctl-js-AAA",
        "metric": "service.http.lvs.443.port",
        "timestamp": ts,
        "step": 60,
        "value": 3,
        "counterType": "MQ,KEYWORD",
        "tags": "",
        "fields": "lvsvip=www.google.com, field2=This metric will overwrite previous value because the (metric - tags - endpoint) are the same."
    },
    {
        "endpoint": "ctl-js-AAA",
        "metric": "service.http.lvs.443.port",
        "timestamp": ts,
        "step": 60,
        "value": 2,
        "counterType": "MQ,KEYWORD",
        "tags": "",
        "fields": "lvsvip=www.google.com, field2=This metric will not overwrite."
    },
]

r = requests.post("http://10.20.30.40:1988/v1/push", data=json.dumps(payload))

print r.text
