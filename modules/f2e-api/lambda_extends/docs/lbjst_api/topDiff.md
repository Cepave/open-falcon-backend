### LBJST: top

* please read common.md first


* API example
* 找出topN的機器
* limit: 回傳上限數字
* orderby: 調整正負排序
  * `desc` `aesc`
* sortby: 指定比較基準數值
  * `Mean` `Max` `Min`

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
    "limit": 3,
    "orderby": "aesc"
  }
}
```
sample response:
```
[
  {
    "Counter": "cpu.idle",
    "DsType": "GAUGE",
    "Endpoint": "hostD",
    "Max": 0,
    "MaxInc": 0,
    "Mean": -3,
    "Min": -5,
    "Step": 60,
    "Values": [
      {
        "Timestamp": 1504239000,
        "Value": 70
      },
      {
        "Timestamp": 1504239060,
        "Value": 75
      },
      {
        "Timestamp": 1504239120,
        "Value": 80
      },
      {
        "Timestamp": 1504239180,
        "Value": 83
      },
      {
        "Timestamp": 1504239240,
        "Value": 83
      }
    ]
  },
  {
    "Counter": "cpu.idle",
    "DsType": "GAUGE",
    "Endpoint": "hostE",
    "Max": 3,
    "MaxInc": 0,
    "Mean": -2,
    "Min": -7,
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
  },
  {
    "Counter": "cpu.idle",
    "DsType": "GAUGE",
    "Endpoint": "hostA",
    "Max": 1,
    "MaxInc": 0,
    "Mean": 0,
    "Min": -2,
    "Step": 60,
    "Values": [
      {
        "Timestamp": 1504238940,
        "Value": 100
      },
      {
        "Timestamp": 1504239000,
        "Value": 99
      },
      {
        "Timestamp": 1504239060,
        "Value": 98
      },
      {
        "Timestamp": 1504239120,
        "Value": 100
      },
      {
        "Timestamp": 1504239180,
        "Value": 100
      }
    ]
  }
]
```
