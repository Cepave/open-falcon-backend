### LBJST: topDiff

* please read common.md first


* API example
* 找出兩點間成長幅度最高/低的metric
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
  "metrices": []string{
    "cpu.idle",
  },
  "func": {
    "limit": 2,
    "orderby": "desc"
  }
}
```
sample reponse:
```
[
  {
    "Counter": "cpu.idle",
    "DsType": "GAUGE",
    "Endpoint": "hostA",
    "Max": 100,
    "Mean": 99,
    "Min": 98,
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
  },
  {
    "Counter": "cpu.idle",
    "DsType": "GAUGE",
    "Endpoint": "hostB",
    "Max": 98,
    "Mean": 88,
    "Min": 80,
    "Step": 60,
    "Values": [
      {
        "Timestamp": 1504239000,
        "Value": 98
      },
      {
        "Timestamp": 1504239060,
        "Value": 89
      },
      {
        "Timestamp": 1504239120,
        "Value": 87
      },
      {
        "Timestamp": 1504239180,
        "Value": 80
      },
      {
        "Timestamp": 1504239240,
        "Value": 85
      }
    ]
  }
]
```
