#!/usr/bin/env python
#!-*- coding:utf8 -*-

import requests
import time
import json

ts = int(time.time())
payload = [
    {
        "endpoint": "test-endpoint",
        "metric": "test-metric",
        "timestamp": ts,
        "step": 60,
        "value": 1,
        "counterType": "GAUGE",
        "tags": "idc=lg,loc=beijing",
        "fields": "strfield1='Mike has a 170mm penis',nonstrfield2=7777"
    },
    {
        "endpoint": "test-endpoint",
        "metric": "mike",
        "timestamp": ts,
        "step": 60,
        "value": 12345,
        "counterType": "MQ,KEYWORD",
        "tags": "idc=lg,loc=beijing",
        "fields": "strfield1='Mike has a 170mm penis',nonstrfield2=7777"
    },
]

r = requests.post("http://10.20.30.40:1988/v1/push", data=json.dumps(payload))

print r.text
