### LBJST: sumAll

* please read common.md first


* API example
* 加總所有的數值 (each timestamp point sum up)
* aliasName: 定義回傳的名稱

sample request:
```
{
  "from":  1504240064,
  "until": 1504250000,
  "endpoints": []string{
    "hostA", "hostB", "hostC", "hostD", "hostE"
  },
  "metrics": []string{
    "cpu.idle",
  },
  "func": {
    "function":  "sumAll",
    "aliasName": "combineAllData"
  }
}
```
sample response:
```
[
  {
    "Counter": "cpu.idle",
    "DsType": "GAUGE",
    "Endpoint": "combineAllData",
    "Step": 60,
    "Values": [
      {
        "Timestamp": 1504238940,
        "Value": 331
      },
      {
        "Timestamp": 1504239000,
        "Value": 326
      },
      {
        "Timestamp": 1504239060,
        "Value": 322
      },
      {
        "Timestamp": 1504239120,
        "Value": 306
      },
      {
        "Timestamp": 1504239180,
        "Value": 323
      }
    ]
  }
]
```
