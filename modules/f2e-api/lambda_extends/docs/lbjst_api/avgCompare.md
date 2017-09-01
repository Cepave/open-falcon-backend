### LBJST: avgCompare

* please read common.md first


* API example
* 找出總值高/低於平均值的機器
* cond: 大於小於
  * ex: > , < , >= , <=, ==

sample request:
```
{
  "from":  1504240064,
  "until": 1504250000,
  "endpoints": []string{
    "hostA", "hostB", "hostC", "hostD", "hostE"
  },
  "metrices": []string{
    "cpu.idle",
  },
  "func": {
    "function":  "avgCompare",
    "cond": "<"
  }
}
```
sample reponse:
```
[
  {
    "Avg": 51.6,
    "Counter": "cpu.idle",
    "DsType": "GAUGE",
    "Endpoint": "hostC",
    "Step": 60,
    "Values": [
      {
        "Timestamp": 1504239000,
        "Value": 60
      },
      {
        "Timestamp": 1504239060,
        "Value": 58
      },
      {
        "Timestamp": 1504239120,
        "Value": 55
      },
      {
        "Timestamp": 1504239180,
        "Value": 40
      },
      {
        "Timestamp": 1504239240,
        "Value": 45
      }
    ]
  },
  {
    "Avg": 4.6,
    "Counter": "cpu.idle",
    "DsType": "GAUGE",
    "Endpoint": "hostE",
    "Step": 60,
    "Values": [
      {
        "Timestamp": 1504239000,
        "Value": 3
      },
      {
        "Timestamp": 1504239060,
        "Value": 5
      },
      {
        "Timestamp": 1504239120,
        "Value": 2
      },
      {
        "Timestamp": 1504239180,
        "Value": 3
      },
      {
        "Timestamp": 1504239240,
        "Value": 10
      }
    ]
  }
]
```
