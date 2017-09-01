### LBJST: Common
* required login session:
  * header: `apiToken: {"name":"testuser92","sig":"48380ba36ad211e79fb3001500c6ca5a"}`

* common api example:
* from: 開始時間
* until: 結束時間
* endpoints: 機器列表
* metrices: 監控項 (with tag)
* func: LBJST參數
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
